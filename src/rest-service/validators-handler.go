package rest

import (
	"errors"
	"fmt"
	"net/http"

	srv "github.com/SeaOfWisdom/sow_library/src/service"
	"github.com/SeaOfWisdom/sow_library/src/service/storage"

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
// @Param 		 Authorization header string true "Bearer {JWT token}"
// @Success      200  {object}   AuthResp
// @Failure      400  {object}  ErrorMsg
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
// @Router       /validator_info [get]
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
// @Param 		 Authorization header string true "Bearer {JWT token}"
// @Success      200  {object}   AuthResp
// @Failure      400  {object}  ErrorMsg
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
// @Router       /work_review [get]
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
