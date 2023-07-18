package rest

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	srv "github.com/SeaOfWisdom/sow_library/src/service"
	"github.com/SeaOfWisdom/sow_library/src/service/storage"
	ocr "github.com/SeaOfWisdom/sow_proto/ocr-srv"

	"github.com/gorilla/mux"
)

/*//////////////////////////

//////// VALIDATOR CRUD ////////

////////*/ ////////////////

// HandleBecomeValidator BecomeValidator godoc
// @Summary      Become a validator
// @Description  Become a validator
// @Tags         Validators
// @Accept       json
// @Produce      json
// @Param        account body BecomeValidatorRequest true "become validator"
// @Success      200  {object}   AuthResp
// @Failure      400  {object}  ErrorMsg
// @Security Bearer
// @Router       /become_validator [post]
func (rs *RestSrv) HandleBecomeValidator(w http.ResponseWriter, r *http.Request) {
	web3Address, err := rs.getWeb3Address(r)
	if err != nil {
		responError(w, http.StatusBadGateway, err.Error())
		return
	}

	request := new(BecomeValidatorRequest)
	if err := rs.getRequest(r.Body, request); err != nil {
		responError(w, http.StatusBadRequest, err.Error())
		return
	}

	// request author
	participant, err := rs.libSrv.BecomeValidator(web3Address, request.EmailAddress, request.Name, request.Surname)
	if err != nil {
		responError(w, http.StatusBadGateway, err.Error())
		return
	}

	// generate a new jwt token for him
	jwt, err := rs.getJWTToken(participant.ID, web3Address, participant.Language, int64(participant.Role))
	if err != nil {
		responError(w, http.StatusBadGateway, "something went wrong, our apologies")
		return
	}

	responJSON(w, http.StatusOK, AuthResp{Token: jwt.Token, Role: storage.ValidatorRole})
}

// HandleValidatorInfo ValidatorInfo godoc
// @Summary      Validator info
// @Description  Validator full info
// @Tags         Validators
// @Accept       json
// @Produce      json
// @Param        web3_address   path      string  true  "validator web3 address"
// @Success      200  {object}   storage.ValidatorResponse
// @Failure      400  {object}  ErrorMsg
// @Router       /validator_info/{web3_address} [get]
func (rs *RestSrv) HandleValidatorInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	validatorAddress, ok := vars["web3_address"]
	if !ok {
		responError(w, http.StatusBadRequest, "null request param")
		return
	}
	rs.logger.Info(fmt.Sprintf("request validator address: %s", validatorAddress))

	// request validator
	validatorResp, err := rs.libSrv.GetValidator(validatorAddress)
	if err != nil {
		responError(w, http.StatusBadGateway, err.Error())
		return
	}

	responJSON(w, http.StatusOK, validatorResp)
}

// HandleUpdateValidator UpdateValidator godoc
// @Summary      Update validator
// @Description  Update validator info
// @Tags         Validators
// @Accept       json
// @Produce      json
// @Param        account body UpdateValidatorRequest true "update validator info"
// @Success      200  {object}   AuthResp
// @Failure      400  {object}  ErrorMsg
// @Security Bearer
// @Router       /update_validator_info [post]
func (rs *RestSrv) HandleUpdateValidator(w http.ResponseWriter, r *http.Request) {
	web3Address, err := rs.getWeb3Address(r)
	if err != nil {
		responError(w, http.StatusBadGateway, err.Error())
		return
	}

	request := new(UpdateValidatorRequest)
	if err := rs.getRequest(r.Body, request); err != nil {
		responError(w, http.StatusBadRequest, err.Error())
		return
	}

	participant, err := rs.libSrv.UpdateValidator(
		web3Address,
		request.EmailAddress,
		request.Name,
		request.Middlename,
		request.Surname,
		request.Orcid,
		request.Language,
		request.Sciences,
	)
	if err != nil {
		responError(w, http.StatusBadRequest, err.Error())
		return
	}

	// generate a new jwt token for him
	jwt, err := rs.getJWTToken(participant.ID, participant.Web3Address, participant.Language, int64(participant.Role))
	if err != nil {
		responError(w, http.StatusBadGateway, "something went wrong, our apologies")
		return
	}

	responJSON(w, http.StatusOK, AuthResp{Token: jwt.Token, Role: participant.Role, NickName: participant.NickName})
}

// HandleUploadValidatorDocs UploadValidatorDocs godoc
// @Summary      Upload validator documents
// @Description  Uploading documents confirming competencies of validator
// @Tags         Validators
// @Accept       json
// @Produce      json
// @Success      200  {object}  SuccessMsg
// @Failure      400  {object}  ErrorMsg
// @Security Bearer
// @Router       /validator_info/upload_docs [post]
func (rs *RestSrv) HandleUploadValidatorDocs(w http.ResponseWriter, r *http.Request) {
	web3Address, err := rs.getWeb3Address(r)
	if err != nil {
		responError(w, http.StatusBadGateway, err.Error())

		return
	}

	participant, err := rs.libSrv.GetParticipantByWeb3Address(web3Address)
	if err != nil {
		responError(w, http.StatusNotFound, err.Error())

		return
	}

	if participant.Role != storage.ValidatorRole {
		responError(w, http.StatusForbidden, "your are not a validator")

		return
	}

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)

	docType := r.FormValue("type")
	if docType == "" {
		responError(w, http.StatusBadRequest, "undefined doc_type param")

		return
	}

	rs.logger.Info(fmt.Sprintf("doc type: %s", docType))

	file, _, err := r.FormFile("doc")
	if err != nil {
		responError(w, http.StatusBadRequest, err.Error())

		return
	}
	defer file.Close()

	// read all of the contents of our uploaded file into a
	// byte array

	/// !!!! TODO !!!!
	imageBytes, err := ioutil.ReadAll(file)
	if err != nil {
		responError(w, http.StatusInternalServerError, "error read file")

		return
	}

	imageResp, err := rs.ocrSrv.ExtractText(r.Context(), &ocr.ExtractTextRequest{
		Image: imageBytes,
	})
	if err != nil {
		err = fmt.Errorf("while extract text via ocr service, err: %v", err)
		rs.logger.Error(err.Error())
		responError(w, http.StatusInternalServerError, err.Error())

		return
	}

	fmt.Println(imageResp)

	responJSON(w, http.StatusOK, SuccessMsg{Msg: "OK"})
}

/*//////////////////////////
///// WORK Evaluation /////
////////////////////////*/

// HandleGetWorkReviewByWorkID GetWorkReviewByWorkID godoc
// @Summary      Work reviews
// @Description  Get all work reviews by work id
// @Tags         Work review
// @Accept       json
// @Produce      json
// @Param        work_id   path      string  true  "work id"
// @Success      200  {object}   storage.WorkReview
// @Failure      400  {object}  ErrorMsg
// @Security Bearer
// @Router       /work_review/{work_id} [get]
func (rs *RestSrv) HandleGetWorkReviewByWorkID(w http.ResponseWriter, r *http.Request) {
	web3Address, err := rs.getWeb3Address(r)
	if err != nil {
		responError(w, http.StatusBadGateway, err.Error())
		return
	}

	vars := mux.Vars(r)
	workId, ok := vars["work_id"]
	if !ok {
		responError(w, http.StatusBadRequest, "null request param")
		return
	}
	rs.logger.Info(fmt.Sprintf("request work id: %s", workId))

	review, err := rs.libSrv.GetValidatorWorkReviewByWorkID(
		web3Address,
		workId,
	)
	if err != nil {
		responError(w, http.StatusBadRequest, err.Error())
		return
	}

	responJSON(w, http.StatusOK, review)
}

// HandleEvaluateWork EvaluateWork godoc
// @Summary      Evaluate work
// @Description  Evaluate work by validator
// @Tags         Work review
// @Accept       json
// @Produce      json
// @Param        account body WorkReviewRequest true "work review"
// @Success      200  {object}   storage.WorkReview
// @Failure      400  {object}  ErrorMsg
// @Security Bearer
// @Router       /update_review [post]
func (rs *RestSrv) HandleEvaluateWork(w http.ResponseWriter, r *http.Request) {
	web3Address, err := rs.getWeb3Address(r)
	if err != nil {
		responError(w, http.StatusBadGateway, err.Error())
		return
	}

	request := new(WorkReviewRequest)
	if err := rs.getRequest(r.Body, request); err != nil {
		responError(w, http.StatusBadRequest, err.Error())
		return
	}

	review, err := rs.libSrv.CreateOrUpdateWorkReview(
		r.Context(),
		web3Address,
		request.Review,
	)
	if err != nil {
		responError(w, http.StatusBadRequest, err.Error())

		return
	}

	responJSON(w, http.StatusOK, review)
}

// HandleSubmitWorkReview SubmitWorkReview godoc
// @Summary      Submit review
// @Description  Submit work review by validator
// @Tags         Work review
// @Accept       json
// @Produce      json
// @Param        work_id   path      string  true  "work id"
// @Param        status   path      string  true "review status" Enums(WORK_REVIEW_SUBMITTED, WORK_REVIEW_SKIPPED, WORK_REVIEW_REJECTED, WORK_REVIEW_DECLINED, WORK_REVIEW_ACCEPTED)
// @Success      200  {object}   SuccessMsg
// @Failure      400  {object}  ErrorMsg
// @Security Bearer
// @Router       /submit_work_review [post]
func (rs *RestSrv) HandleSubmitWorkReview(w http.ResponseWriter, r *http.Request) {
	web3Address, err := rs.getWeb3Address(r)
	if err != nil {
		responError(w, http.StatusBadGateway, err.Error())

		return
	}

	vars := mux.Vars(r)
	workID, ok := vars["work_id"]
	if !ok {
		responError(w, http.StatusBadRequest, "null work_id path param")

		return
	}

	statusAsString := vars["status"]

	rs.logger.Info(fmt.Sprintf("request work id: %s", workID))

	err = rs.libSrv.SubmitWorkReview(r.Context(), web3Address, workID, storage.StringToReviewStatus(statusAsString))
	if err != nil {
		if errors.Is(err, srv.ErrNoReviews) {
			responError(w, http.StatusNotFound, err.Error())

			return
		}

		responError(w, http.StatusInternalServerError, err.Error())

		return
	}

	responJSON(w, http.StatusOK, SuccessMsg{Msg: "OK"})
}
