package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	_ "github.com/SeaOfWisdom/sow_library/docs" // docs is generated by Swag CLI, you have to import it.
	jwt "github.com/SeaOfWisdom/sow_proto/jwt-srv"
	ocr "github.com/SeaOfWisdom/sow_proto/ocr-srv"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/SeaOfWisdom/sow_library/src/config"
	srv "github.com/SeaOfWisdom/sow_library/src/service"
	"github.com/SeaOfWisdom/sow_library/src/service/storage"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type IRequest interface {
	Validate() error
}

var authURIs map[string]storage.ParticipantRole

type RestSrv struct {
	logger *zap.Logger

	router *mux.Router
	server *http.Server

	libSrv *srv.LibrarySrv

	/* grpc services */
	jwtSrv jwt.JwtServiceClient
	ocrSrv ocr.OCRClient
}

func NewRestSrv(
	cfg *config.Config,
	libSrv *srv.LibrarySrv,
	jwtSrv jwt.JwtServiceClient,
	ocrSrv ocr.OCRClient,
) *RestSrv {
	instance := &RestSrv{
		logger: zap.NewExample(),
		router: mux.NewRouter(),
		server: &http.Server{
			Addr: cfg.RestAddress,
		},
		libSrv: libSrv,
		jwtSrv: jwtSrv,
		ocrSrv: ocrSrv,
	}
	// set router && create ethereum service
	instance.router = mux.NewRouter()
	authURIs = instance.setRouters()

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})
	instance.server.Handler = handlers.CORS(headers, methods, origins)(instance.router)

	return instance
}

// Get wraps the router for GET method
func (rs *RestSrv) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	rs.router.HandleFunc(path, f).Methods("GET")
}

// Get wraps the router for POST method
func (rs *RestSrv) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	rs.router.HandleFunc(path, f).Methods("POST")
}

// Get wraps the router for PUT method
func (rs *RestSrv) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	rs.router.HandleFunc(path, f).Methods("PUT")
}

func (rs *RestSrv) setRouters() map[string]storage.ParticipantRole {
	// Participants
	rs.Post("/new_participant", rs.HandleNewParticipant)
	rs.Get("/auth/{web3_address}", rs.HandleAuth)
	rs.Get("/get_basic_info", rs.HandleGetBasicInfo)
	rs.Get("/if_participant_exists/{web3_address}", rs.HandleIfParticipantExists)

	rs.Post("/update_basic_info", rs.HandleUpdataBasicPacticipant)

	// rs.Post("/remove_participant/{nick_name}", rs.HandleAddPacticipant)

	// Author methods

	rs.Get("/author_data", rs.HandleBecomeAuthorData)
	// returns the data related the registration process
	rs.Post("/become_author", rs.HandleBecomeAuthor)
	rs.Get("/author_info/{web3_address}", rs.HandleAuthorInfo)
	rs.Post("/invite_co_author", rs.HandleInviteCoAuthor)
	rs.Post("/update_author_info", rs.HandleUpdateAuthor)

	// Validator methods TODO
	rs.Post("/become_validator", rs.HandleBecomeValidator)
	rs.Post("/update_validator_info", rs.HandleUpdateValidator)
	// TODO
	rs.Get("/validator_info/{web3_address}", rs.HandleValidatorInfo)

	// Validator work review
	rs.Get("/work_review/{work_id}", rs.HandleGetWorkReviewByWorkID)
	rs.Get("/work_reviews/{work_id}", rs.HandleGetWorkReviews)
	rs.Post("/update_review", rs.HandleEvaluateWork)
	rs.Post("/submit_work_review/{work_id}/{status}", rs.HandleSubmitWorkReview)

	// DOCs
	rs.Put("/upload_doc/{doc_type}", rs.HandlerUploadDoc)

	// Works
	rs.Get("/works", rs.HandleAllWorks)
	// TODO
	rs.Get("/works/{work_id}", rs.HandleWorkByID)
	rs.Get("/works/author/{web3_address}", rs.HandleAuthorWorks)

	rs.Get("/works_by_key_words/{key_words}", rs.HandleWorkByKeyWords)

	rs.Get("/purchase_work/{work_id}", rs.HandlePurchaseWork)
	rs.Get("/purchased_works", rs.HandlePurchasedWorks)

	// Bookmarks
	rs.Post("/add_bookmark/{work_id}", rs.HandleAddInBookmarks)
	rs.Post("/remove_bookmark/{work_id}", rs.HandleRemoveFromBookmarks)
	rs.Get("/bookmarks", rs.HandleGetBookmarks)

	rs.Get("/work_data", rs.HandlePublishWorkData)
	rs.Post("/publish_work", rs.HandlePublishWork)

	// only admin role
	rs.Get("/pending_works", rs.HandlePendingWorks)
	rs.Post("/approve_work/{work_id}", rs.HandleApproveWork)
	rs.Post("/remove_work/{work_id}", rs.HandleRemoveWork)

	rs.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		// httpSwagger.URL("http://0.0.0.0:8005/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	rs.router.Use(rs.jwtMiddleware)
	return map[string]storage.ParticipantRole{
		// Participants
		"get_basic_info":    storage.ReaderRole,
		"update_basic_info": storage.ReaderRole,

		//	"/remove_participant": storage.AdminRole,
		"become_author": storage.ReaderRole,

		"update_author_info": storage.AuthorRole,

		"become_validator":   storage.ReaderRole,
		"update_review":      storage.ValidatorRole,
		"work_review":        storage.ValidatorRole,
		"submit_work_review": storage.ValidatorRole,

		// Works
		"publish_work":  storage.AuthorRole,
		"pending_works": storage.AdminRole,
		"approve_work":  storage.AdminRole,
		"remove_work":   storage.AdminRole,

		"purchase_work":   storage.ReaderRole,
		"purchased_works": storage.ReaderRole,

		// Bookmarks
		"add_bookmark":    storage.ReaderRole,
		"remove_bookmark": storage.ReaderRole,
		"bookmarks":       storage.ReaderRole,
	}
}

// Run the app on it's router
func (rs *RestSrv) Start() {
	go func() {
		if err := rs.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("failed to start rest-server, err: %v", err))
		}
	}()
	rs.logger.Info(fmt.Sprintf("rest-server was started on %s", rs.server.Addr))
}

// respondJSON makes the response with payload as json format
func responJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

type SuccessMsg struct {
	Msg string `json:"status" example:"OK"`
}

type ErrorMsg struct {
	Msg string `json:"error" example:"null request param"`
}

// respondError makes the error response with payload as json format
func responError(w http.ResponseWriter, code int, message string) {
	responJSON(w, code, ErrorMsg{Msg: message})
}

func (rs *RestSrv) getRequest(reader io.Reader, request interface{}) error {
	reqBody, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Printf("server: could not read request body: %v\n", err)
		return err
	}
	if err := json.Unmarshal(reqBody, &request); err != nil {
		fmt.Printf("server: could not read request body: %v\n", err)
		return err
	}

	if err := request.(IRequest).Validate(); err != nil {
		fmt.Printf("server: while validating, err %v\n", err)
		return err
	}
	rs.logger.Info(fmt.Sprintf("request: %v", request))
	return nil
}

func (rs *RestSrv) getJWTToken(sub, web3Address, language string, role int64) (*jwt.Token, error) {
	// generate a jwt token for him
	req := &jwt.TokenBody{Iss: srv.ServiceName, Sub: sub, Role: role, Web3Address: web3Address}
	resp, err := rs.jwtSrv.GenerateJWT(context.Background(), req)
	if err != nil {
		rs.logger.Error(fmt.Sprintf("generate jwt, err: %v", err))
		return nil, err
	}
	return resp, nil
}

func (rs *RestSrv) Stop() {
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rs.server.Shutdown(ctxShutDown); err != nil {
		rs.logger.Fatal(fmt.Sprintf("failed to stop rest-server, err: %v", err))
	}
}