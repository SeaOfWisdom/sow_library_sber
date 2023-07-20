package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// HandleFaucet  Faucet godoc
// @Summary      Faucet SOW tokens
// @Description  Mints 50 SOW tokens to web3_address
// @Tags         Faucet
// @Accept       json
// @Produce      json
// @Param        web3_address   path      string  true  "participant web3 address"
// @Success      200  {object}   string
// @Failure      400  {object}  ErrorMsg
// @Router       /faucet [get]
func (rs *RestSrv) HandleFaucet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	web3AddrStr, ok := vars["web3_address"]
	if !ok {
		responError(w, http.StatusBadRequest, "null request param")

		return
	}

	rs.logger.Infof("HandleFaucet: request address: %s", web3AddrStr)

	txHash, err := rs.libSrv.Faucet(web3AddrStr)
	if err != nil {
		responError(w, http.StatusInternalServerError, err.Error())

		return
	}

	responJSON(w, http.StatusOK, txHash)
}

// HandleAuth Auth godoc
// @Summary      Auth account
// @Description  Auth account and return JWT token
// @Tags         Authorization
// @Accept       json
// @Produce      json
// @Param        web3_address   path      string  true  "participant web3 address"
// @Success      200  {object}   AuthResp
// @Failure      400  {object}  ErrorMsg
// @Router       /auth/{web3_address} [get]
func (rs *RestSrv) HandleAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	web3AddrStr, ok := vars["web3_address"]
	if !ok {
		responError(w, http.StatusBadRequest, "null request param")

		return
	}

	rs.logger.Info("HandleAuth: request address: %s", web3AddrStr)

	// try to find a participant by the web3 address
	participant, err := rs.libSrv.GetParticipantByWeb3Address(web3AddrStr)
	if err != nil {
		responError(w, http.StatusInternalServerError, err.Error())

		return
	}

	// generate a new jwt token for him
	jwt, err := rs.getJWTToken(ctx, participant.ID, participant.Web3Address, participant.Language, int64(participant.Role))
	if err != nil {
		responError(w, http.StatusInternalServerError, "something went wrong, our apologies")

		return
	}

	responJSON(w, http.StatusOK, AuthResp{Token: jwt.Token, Role: participant.Role, NickName: participant.NickName})
}

// HandleNewParticipant NewParticipant godoc
// @Summary      Become a participant
// @Description  Become a new participant
// @Tags         Participants
// @Accept       json
// @Produce      json
// @Param        account body NewParticipantRequest true "become participant"
// @Success      200  {object}   AuthResp
// @Failure      400  {object}  ErrorMsg
// @Router       /new_participant [post]
func (rs *RestSrv) HandleNewParticipant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request := new(NewParticipantRequest)
	if err := rs.getRequest(r.Body, request); err != nil {
		responError(w, http.StatusBadRequest, err.Error())

		return
	}

	participant, err := rs.libSrv.CreateParticipant(ctx, request.NickName, request.Web3Address)
	if err != nil {
		responError(w, http.StatusInternalServerError, err.Error())

		return
	}

	// generate a new jwt token for him
	jwt, err := rs.getJWTToken(ctx, participant.ID, participant.Web3Address, participant.Language, int64(participant.Role))
	if err != nil {
		responError(w, http.StatusInternalServerError, "something went wrong, our apologies")

		return
	}

	responJSON(w, http.StatusOK, AuthResp{Token: jwt.Token, Role: participant.Role, NickName: request.NickName})
}

// HandleIfParticipantExists IfParticipantExists godoc
// @Summary      Check if participant exists
// @Description  Check participant availability
// @Tags         Participants
// @Accept       json
// @Produce      json
// @Param        web3_address   path      string  true  "participant web3 address"
// @Success      200  {object}   IfParticipantExistsResp
// @Failure      400  {object}  ErrorMsg
// @Router       /if_participant_exists/{web3_address} [get]
func (rs *RestSrv) HandleIfParticipantExists(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	web3AddrStr, ok := vars["web3_address"]
	if !ok {
		responError(w, http.StatusBadRequest, "null request param")

		return
	}
	rs.logger.Info(fmt.Sprintf("request address: %s", web3AddrStr))

	// try to find a participant by the web3 address
	if _, err := rs.libSrv.GetParticipantByWeb3Address(web3AddrStr); err != nil {
		responJSON(w, http.StatusOK, IfParticipantExistsResp{Result: false})

		return
	}

	responJSON(w, http.StatusOK, IfParticipantExistsResp{Result: true})
}

// HandleUpdateBasicParticipant UpdateBasicParticipant godoc
// @Summary     Update participant info
// @Description Update basic participant info
// @Tags		Participants
// @Accept      json
// @Produce     json
// @Param       account body BasicInfoUpdateRequest true "update basic participant info"
// @Success     200  {object}   AuthResp
// @Failure     400  {object}  ErrorMsg
// @Security Bearer
// @Router		/update_basic_info [post]
func (rs *RestSrv) HandleUpdateBasicParticipant(w http.ResponseWriter, r *http.Request) {
	web3Address, err := rs.getWeb3Address(r)
	if err != nil {
		responError(w, http.StatusUnauthorized, err.Error())

		return
	}

	request := new(BasicInfoUpdateRequest)
	if err := rs.getRequest(r.Body, request); err != nil {
		responError(w, http.StatusBadRequest, err.Error())

		return
	}

	participant, err := rs.libSrv.UpdateParticipantNickName(web3Address, request.NickName)
	if err != nil {
		responError(w, http.StatusInternalServerError, err.Error())

		return
	}

	// generate a new jwt token for him
	jwt, err := rs.getJWTToken(r.Context(), participant.ID, web3Address, participant.Language, int64(participant.Role))
	if err != nil {
		responError(w, http.StatusInternalServerError, "something went wrong, our apologies")

		return
	}

	responJSON(w, http.StatusOK, AuthResp{Token: jwt.Token, Role: participant.Role, NickName: request.NickName})
}

// HandleGetBasicInfo GetBasicInfo godoc
// @Summary		Get info
// @Description Get basic info
// @Tags        Participants
// @Accept      json
// @Produce     json
// @Success     200  {object}  BasicInfo
// @Failure     400  {object}  ErrorMsg
// @Security Bearer
// @Router      /get_basic_info [post]
func (rs *RestSrv) HandleGetBasicInfo(w http.ResponseWriter, r *http.Request) {
	web3Address, err := rs.getWeb3Address(r)
	if err != nil {
		responError(w, http.StatusUnauthorized, fmt.Sprintf("while getting the decoding the jwt token, err: %v", err))

		return
	}

	// get nickname
	participant, err := rs.libSrv.GetParticipantByWeb3Address(web3Address)
	if err != nil {
		responError(w, http.StatusInternalServerError, err.Error())

		return
	}

	responJSON(w, http.StatusOK, BasicInfo{NickName: participant.NickName, Role: participant.Role})
}
