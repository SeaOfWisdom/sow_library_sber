package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// HandleWorkByID Work by id
// @Summary      Work by id
// @Description  Get work by id
// @Tags         Works
// @Accept       json
// @Produce      json
// @Param        work_id   path      string  true  "work id"
// @Success 	200 {object} []storage.WorkResponse
// @Failure      400  {object}  ErrorMsg
// @Failure      404  {object}  ErrorMsg
// @Security Bearer
// @Router       /works/{work_id} [get]
func (rs *RestSrv) HandleWorkByID(w http.ResponseWriter, r *http.Request) {
	web3Address, err := rs.getWeb3Address(r)
	if err != nil {
		responError(w, http.StatusUnauthorized, fmt.Sprintf("while getting the decoding the jwt token, err: %v", err))

		return
	}

	vars := mux.Vars(r)
	workID, ok := vars["work_id"]
	if !ok {
		responError(w, http.StatusBadRequest, "null request param")

		return
	}

	rs.logger.Infof("request work id: %s", workID)

	work, err := rs.libSrv.GetWorkByID(r.Context(), web3Address, workID)
	if err != nil {
		responError(w, http.StatusInternalServerError, err.Error())

		return
	}

	if work == nil {
		responError(w, http.StatusNotFound, "work was not found")

		return
	}

	responJSON(w, http.StatusOK, work)
}

// HandleAllWorks AllWorks godoc
// @Summary      Get all works
// @Description  Get all works depends on role
// @Tags         Works
// @Accept       json
// @Produce      json
// @Success 	200 {object} []storage.WorkResponse
// @Failure      400  {object}  ErrorMsg
// @Router       /works [get]
func (rs *RestSrv) HandleAllWorks(w http.ResponseWriter, r *http.Request) {
	web3Address, err := rs.getWeb3Address(r)
	if err != nil {
		if errors.Is(err, ErrNoToken) {
			web3Address = ""
		} else {
			responError(w, http.StatusUnauthorized, fmt.Sprintf("while getting the decoding the jwt token, err: %v", err))

			return
		}
	}

	works, err := rs.libSrv.GetAllWorks(r.Context(), web3Address)
	if err != nil {
		responError(w, http.StatusBadRequest, err.Error())

		return
	}

	responJSON(w, http.StatusOK, works)
}

// HandleWorkByKeyWords WorkByKeyWords godoc
// @Summary      Get all works by words
// @Description  Search works by particular words
// @Tags         Works
// @Accept       json
// @Produce      json
// @Param        key_words   path      string  true  "words to search for"
// @Success 	200 {object} []storage.WorkResponse
// @Failure      400  {object}  ErrorMsg
// @Security Bearer
// @Router       /works_by_key_words [get]
func (rs *RestSrv) HandleWorkByKeyWords(w http.ResponseWriter, r *http.Request) {
	web3Address, err := rs.getWeb3Address(r)
	if err != nil {
		if strings.Contains(err.Error(), "missing auth token") ||
			strings.Contains(err.Error(), "invalid/Malformed auth token") {
			web3Address = ""
		} else {
			responError(w, http.StatusUnauthorized, fmt.Sprintf("while getting the decoding the jwt token, err: %v", err))

			return
		}
	}

	vars := mux.Vars(r)
	keyWordsReq, ok := vars["key_words"]
	if !ok {
		responError(w, http.StatusBadRequest, "null request param")

		return
	}
	rs.logger.Infof("request key words: %s", keyWordsReq)

	keyWords := strings.Split(keyWordsReq, ",")
	if len(keyWords) == 0 {
		responError(w, http.StatusBadRequest, "null request param")

		return
	}

	works, err := rs.libSrv.GetWorksByKeyWords(r.Context(), web3Address, keyWords)
	if err != nil {
		responError(w, http.StatusInternalServerError, err.Error())

		return
	}
	if works == nil {
		responError(w, http.StatusOK, "there are no works in the library")

		return
	}

	responJSON(w, http.StatusOK, works)
}

// HandlePublishWorkData PublishWorkData godoc
// @Summary      Mock work data
// @Description  Mock work data
// @Tags         Publish work
// @Accept       json
// @Produce      json
// @Success 	200 {object} PublishWorkDataResp
// @Router       /work_data [get]
func (rs *RestSrv) HandlePublishWorkData(w http.ResponseWriter, r *http.Request) {
	responJSON(w, http.StatusOK, PublishWorkDataResp{Tags: []string{
		"подводный спорт", "моноласт", "плавание в ластах", "скоростное плавание", "апноэ"}})
}

// HandlePublishWork PublishWork godoc
// @Summary      Publish a new work
// @Description  Publish a new work
// @Tags         Publish work
// @Accept       json
// @Produce      json
// @Param		Authorization header string true "Bearer {JWT token}"
// @Param		Work body storage.WorkResponse true "work and author info"
// @Success 	200 {object} storage.WorkResponse
// @Failure      400  {object}  ErrorMsg
// @Failure      401  {object}  ErrorMsg
// @Security Bearer
// @Router       /publish_work [post]
func (rs *RestSrv) HandlePublishWork(w http.ResponseWriter, r *http.Request) {
	// get address from the JWT token
	web3Address, err := rs.getWeb3Address(r)
	if err != nil {
		responError(w, http.StatusUnauthorized, fmt.Sprintf("while getting the decoding the jwt token, err: %v", err))

		return
	}

	request := new(WorkReq)
	if err := rs.getRequest(r.Body, request); err != nil {
		responError(w, http.StatusInternalServerError, err.Error())

		return
	}

	// TODO
	workResp, _, err := rs.libSrv.PublishWork(r.Context(), web3Address, request.Work)
	if err != nil {
		responError(w, http.StatusInternalServerError, err.Error())

		return
	}

	responJSON(w, http.StatusOK, workResp)
}

// HandleAuthorWorks AuthorWorks godoc
// @Summary      List author`s works
// @Description  Get list author`s works
// @Tags         Works
// @Accept       json
// @Produce      json
// @Param        web3_address   path      string  true  "author web3 address"
// @Success 	200 {object} []storage.WorkResponse
// @Failure      400  {object}  ErrorMsg
// @Security Bearer
// @Router       /works/author/{web3_address} [get]
func (rs *RestSrv) HandleAuthorWorks(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// defer ctx.Done()
	// role := ctx.Value("role").(storage.ParticipantRole)
	// get address from the JWT token
	web3Address, err := rs.getWeb3Address(r)
	if err != nil {
		responError(w, http.StatusUnauthorized, err.Error())

		return
	}

	vars := mux.Vars(r)
	we3Address, ok := vars["web3_address"]
	if !ok {
		responError(w, http.StatusBadRequest, "null request param")

		return
	}

	rs.logger.Infof("request author's (%s) works", we3Address)

	works, err := rs.libSrv.GetWorksByAuthorAddress(r.Context(), web3Address, we3Address)
	if err != nil {
		responError(w, http.StatusBadRequest, err.Error())

		return
	}
	if works == nil {
		responJSON(w, http.StatusOK, "there are no pending works")

		return
	}

	responJSON(w, http.StatusOK, works)
}

/*//////////////////////////////
///// Purchasing works /////
////////////////////////////*/

// HandlePurchaseWork PurchaseWork godoc
// @Summary      Purchase work
// @Description  Purchase particular work
// @Tags         Purchasing works
// @Accept       json
// @Produce      json
// @Param        work_id   path      string  true  "work id to purchase"
// @Success 	200 {object} SuccessMsg
// @Failure      400  {object}  ErrorMsg
// @Failure      401  {object}  ErrorMsg
// @Security Bearer
// @Router       /purchase_work/{work_id} [post]
func (rs *RestSrv) HandlePurchaseWork(w http.ResponseWriter, r *http.Request) {
	// get address from the JWT token
	web3Address, err := rs.getWeb3Address(r)
	if err != nil {
		responError(w, http.StatusUnauthorized, fmt.Sprintf("while getting the decoding the jwt token, err: %v", err))

		return
	}

	vars := mux.Vars(r)
	workID, ok := vars["work_id"]
	if !ok {
		responError(w, http.StatusBadRequest, "null request param")

		return
	}
	rs.logger.Info(fmt.Sprintf("request work id: %s", workID))

	if err := rs.libSrv.PurchaseWork(r.Context(), web3Address, workID, false); err != nil {
		responError(w, http.StatusBadRequest, err.Error())

		return
	}

	responJSON(w, http.StatusOK, SuccessMsg{Msg: "OK"})
}

// HandlePurchasedWorks PurchasedWorks godoc
// @Summary      Purchased works
// @Description  Get purchased works
// @Tags         Purchasing works
// @Accept       json
// @Produce      json
// @Success 	200 {object} []storage.WorkResponse
// @Failure      400  {object}  ErrorMsg
// @Failure      401  {object}  ErrorMsg
// @Security Bearer
// @Router       /purchased_works [get]
func (rs *RestSrv) HandlePurchasedWorks(w http.ResponseWriter, r *http.Request) {
	// get address from the JWT token
	web3Address, err := rs.getWeb3Address(r)
	if err != nil {
		responError(w, http.StatusUnauthorized, fmt.Sprintf("while getting the decoding the jwt token, err: %v", err))

		return
	}

	works, err := rs.libSrv.PurchasedWorks(r.Context(), web3Address)
	if err != nil {
		responError(w, http.StatusBadRequest, err.Error())

		return
	}

	responJSON(w, http.StatusOK, works)
}

/*////////////////////
///// Bookmarks /////
//////////////////*/

// HandleAddInBookmarks AddInBookmarks godoc
// @Summary      Add a new bookmark
// @Description  Add a new bookmark to work
// @Tags		 Bookmarks
// @Accept       json
// @Produce      json
// @Param        work_id   path      string  true  "work id to add into bookmarks"
// @Success 	200 {object} []storage.WorkResponse
// @Failure      400  {object}  ErrorMsg
// @Failure      401  {object}  ErrorMsg
// @Security Bearer
// @Router       /add_bookmark/{work_id} [get]
func (rs *RestSrv) HandleAddInBookmarks(w http.ResponseWriter, r *http.Request) {
	// get the reader's web3 address from the JWT token
	web3Address, err := rs.getWeb3Address(r)
	if err != nil {
		rs.logger.Errorf("HandleAddInBookmarks: %v", err)
		responError(w, http.StatusUnauthorized, err.Error())

		return
	}
	// get the work id from the request params
	vars := mux.Vars(r)
	workId, ok := vars["work_id"]
	if !ok {
		rs.logger.Errorf("HandleAddInBookmarks: %v", err)
		responError(w, http.StatusBadRequest, "null request param")

		return
	}
	rs.logger.Info(fmt.Sprintf("request work id: %s", workId))

	if err := rs.libSrv.CreateBookmark(web3Address, workId); err != nil {
		rs.logger.Errorf("HandleAddInBookmarks: %v", err)
		responError(w, http.StatusInternalServerError, fmt.Sprintf("something went wrong while creating a new bookmark, err: %v", err))

		return
	}

	responJSON(w, http.StatusOK, SuccessMsg{Msg: "OK"})
}

// HandleGetBookmarks GetBookmarks godoc
// @Summary      Get bookmarks
// @Description  Get bookmarks
// @Tags		 Bookmarks
// @Accept       json
// @Produce      json
// @Success 	200 {object} []storage.WorkResponse
// @Failure      400  {object}  ErrorMsg
// @Failure      401  {object}  ErrorMsg
// @Security Bearer
// @Router       /bookmarks [get]
func (rs *RestSrv) HandleGetBookmarks(w http.ResponseWriter, r *http.Request) {
	// get the reader's web3 address from the JWT token
	web3Address, err := rs.getWeb3Address(r)
	if err != nil {
		rs.logger.Errorf("HandleGetBookmarks: %v", err)
		responError(w, http.StatusUnauthorized, err.Error())

		return
	}

	bookmarks, err := rs.libSrv.GetBookmarksOf(r.Context(), web3Address)
	if err != nil {
		rs.logger.Errorf("HandleGetBookmarks: %v", err)
		responError(w, http.StatusInternalServerError, err.Error())

		return
	}

	responJSON(w, http.StatusOK, bookmarks)
}

// HandleRemoveFromBookmarks RemoveFromBookmarks godoc
// @Summary      Remove bookmarks
// @Description  Remove bookmarks
// @Tags		 Bookmarks
// @Accept       json
// @Produce      json
// @Param        web3_address   path      string  true  "author web3 address"
// @Param		Authorization header string true "Bearer {JWT token}"
// @Success 	200 {object} SuccessMsg
// @Failure      400  {object}  ErrorMsg
// @Failure      401  {object}  ErrorMsg
// @Security Bearer
// @Router       /remove_bookmark [post]
func (rs *RestSrv) HandleRemoveFromBookmarks(w http.ResponseWriter, r *http.Request) {
	// get the reader's web3 address from the JWT token
	web3Address, err := rs.getWeb3Address(r)
	if err != nil {
		rs.logger.Errorf("HandleRemoveFromBookmarks: %v", err)
		responError(w, http.StatusUnauthorized, err.Error())

		return
	}
	// get the work id from the request params
	vars := mux.Vars(r)
	workId, ok := vars["work_id"]
	if !ok {
		rs.logger.Errorf("HandleRemoveFromBookmarks: %v", err)
		responError(w, http.StatusBadRequest, "null request param")

		return
	}
	rs.logger.Info("HandleRemoveFromBookmarks: request work id: %s", workId)

	if err := rs.libSrv.RemoveBookmark(web3Address, workId); err != nil {
		rs.logger.Errorf("HandleRemoveFromBookmarks: %v", err)
		responError(w, http.StatusInternalServerError, err.Error())

		return
	}

	responJSON(w, http.StatusOK, SuccessMsg{Msg: "OK"})
}

/// ----- ----- Admin methods ----- -----

func (rs *RestSrv) HandleApproveWork(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	workId, ok := vars["work_id"]
	if !ok {
		responError(w, http.StatusBadRequest, "null request param")

		return
	}
	rs.logger.Infof("HandleApproveWork: request work id: %s", workId)

	if err := rs.libSrv.ApproveWork(r.Context(), workId); err != nil {
		responError(w, http.StatusInternalServerError, err.Error())

		return
	}
	responJSON(w, http.StatusOK, "ok") // TODO
}

func (rs *RestSrv) HandlePendingWorks(w http.ResponseWriter, r *http.Request) {
	works, err := rs.libSrv.GetPendingWorks(r.Context())
	if err != nil {
		responError(w, http.StatusInternalServerError, err.Error())

		return
	}
	if works == nil {
		responJSON(w, http.StatusOK, "there are no pending works")

		return
	}
	responJSON(w, http.StatusOK, works)
}

func (rs *RestSrv) HandleRemoveWork(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	workId, ok := vars["work_id"]
	if !ok {
		responError(w, http.StatusBadRequest, "null request param")

		return
	}
	rs.logger.Infof("HandleRemoveWork: request work id: %s", workId)

	if err := rs.libSrv.RemoveWork(r.Context(), workId); err != nil {
		responError(w, http.StatusInternalServerError, err.Error())

		return
	}

	responJSON(w, http.StatusOK, "OK") // TODO
}
