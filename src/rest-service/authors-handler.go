package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// HandleBecomeAuthorData
// @Summary Get sciences
// @Description  Get all available sciences
// @Tags         Authors
// @Accept       json
// @Produce      json
// @Success      200  {object}  BecomeAuthorDataResp
// @Router       /author_data [get]
func (rs *RestSrv) HandleBecomeAuthorData(w http.ResponseWriter, r *http.Request) {
	responJSON(w, http.StatusOK, BecomeAuthorDataResp{
		Sciences: []string{"13.00.04-Теория и методика физического воспитания, спортивной тренировки, оздоровительной и адаптивной физической культуры", "13.00.08-Теория и методика профессионального образования"},
	})
}

// HandleBecomeAuthor BecomeAuthor godoc
// @Summary      Become a author
// @Description  Become a author
// @Tags         Authors
// @Accept       json
// @Produce      json
// @Param        account body BecomeAuthorRequest true "update author info"
// @Success      200  {object}  AuthResp
// @Failure      400  {object}  ErrorMsg
// @Security Bearer
// @Router       /become_author [post]
func (rs *RestSrv) HandleBecomeAuthor(w http.ResponseWriter, r *http.Request) {
	web3Address, err := rs.getWeb3Address(r)
	if err != nil {
		responError(w, http.StatusBadGateway, err.Error())
		return
	}

	request := new(BecomeAuthorRequest)
	if err := rs.getRequest(r.Body, request); err != nil {
		responError(w, http.StatusBadRequest, err.Error())
		return
	}

	// request author
	participant, err := rs.libSrv.BecomeAuthor(web3Address, request.EmailAddress, request.Name, request.Surname)
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

	responJSON(w, http.StatusOK, AuthResp{Token: jwt.Token, Role: participant.Role})
}

// HandleInviteCoAuthor
// @Summary      Invite co-author
// @Description  Invite co-author
// @Tags         Authors
// @Accept       json
// @Produce      json
// @Param        account body BecomeAuthorRequest true "update author info"
// @Success      200  {object}  SuccessMsg
// @Failure      400  {object}  ErrorMsg
// @Security Bearer
// @Router       /invite_co_author [post]
func (rs *RestSrv) HandleInviteCoAuthor(w http.ResponseWriter, r *http.Request) {
	request := new(BecomeAuthorRequest)
	if err := rs.getRequest(r.Body, request); err != nil {
		responError(w, http.StatusBadRequest, err.Error())
		return
	}

	// TODO
	responJSON(w, http.StatusOK, SuccessMsg{Msg: "OK"})
}

// HandleAuthorInfo AuthorInfo godoc
// @Summary      Get author info
// @Description  Get author info
// @Tags         Authors
// @Accept       json
// @Produce      json
// @Param        web3_address   path      string  true  "author web3 address"
// @Success      200  {object}  storage.AuthorResponse
// @Failure      400  {object}  ErrorMsg
// @Router       /author_info/{web3_address} [get]
func (rs *RestSrv) HandleAuthorInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authorAddress, ok := vars["web3_address"]
	if !ok {
		responError(w, http.StatusBadRequest, "null request param")
		return
	}
	rs.logger.Info(fmt.Sprintf("request author address: %s", authorAddress))

	// request author
	authorResp, err := rs.libSrv.GetAuthor(authorAddress)
	if err != nil {
		responError(w, http.StatusBadGateway, err.Error())
		return
	}

	responJSON(w, http.StatusOK, authorResp)
}

// HandleUpdateAuthor
// @Summary      Update participant info
// @Description  Update participant basic info
// @Tags         Authors
// @Accept       json
// @Produce      json
// @Param        account body UpdateAuthorRequest true "update author info"
// @Success      200  {object}   AuthResp
// @Failure      400  {object}  ErrorMsg
// @Security Bearer
// @Router       /update_author_info [post]
func (rs *RestSrv) HandleUpdateAuthor(w http.ResponseWriter, r *http.Request) {
	web3Address, err := rs.getWeb3Address(r)
	if err != nil {
		responError(w, http.StatusBadGateway, err.Error())
		return
	}

	request := new(UpdateAuthorRequest)
	if err := rs.getRequest(r.Body, request); err != nil {
		responError(w, http.StatusBadRequest, err.Error())
		return
	}

	participant, err := rs.libSrv.UpdateAuthorInfo(
		web3Address,
		request.EmailAddress,
		request.Name,
		request.Middlename,
		request.Surname,
		request.Orcid,
		request.ScholarShipProfile,
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

// HandleGetWorkReviews godoc
// @Summary      Get work reviews
// @Description Get work reviews by work_id
// @Tags         Work review
// @Accept       json
// @Produce      json
// @Param        work_id   path      string  true  "work id"
// @Success      200  {object}   []storage.WorkReview
// @Failure      400  {object}  ErrorMsg
// @Security Bearer
// @Router       /work_reviews/{work_id} [get]
func (rs *RestSrv) HandleGetWorkReviews(w http.ResponseWriter, r *http.Request) {
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

	review, err := rs.libSrv.GetWorkReviewsByWorkID(
		web3Address,
		workId,
	)
	if err != nil {
		responError(w, http.StatusBadRequest, err.Error())
		return
	}

	responJSON(w, http.StatusOK, review)
}
