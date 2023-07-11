package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// HandleAuth Auth godoc
// @Summary      Auth account
// @Description  Auth account and return JWT token
// @Tags         Authorization
// @Accept       json
// @Produce      json
// @Param        web3_address   path      string  true  "participant web3 address"
// @Success      200  {object}   AuthResp
// @Failure      400  {object}  ErrorMsg
// @Router       /auth [get]
func (rs *RestSrv) HandleAuth(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	web3AddrStr, ok := vars["web3_address"]
	if !ok {
		responError(w, http.StatusBadRequest, "null request param")
		return
	}
	rs.logger.Info(fmt.Sprintf("request address: %s", web3AddrStr))

	// try to find a participant by the web3 address
	participant, err := rs.libSrv.GetParticipantByWeb3Address(web3AddrStr)
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
	request := new(NewParticipantRequest)
	if err := rs.getRequest(r.Body, request); err != nil {
		responError(w, http.StatusBadRequest, err.Error())
		return
	}

	participant, err := rs.libSrv.CreateParticipant(request.NickName, request.Web3Address)
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

	responJSON(w, http.StatusOK, AuthResp{Token: jwt.Token, Role: participant.Role, NickName: request.NickName})
}

// IfParticipantExists godoc
// @Summary      TODO
// @Description  TODO
// @Tags         Participants
// @Accept       json
// @Produce      json
// @Param        web3_address   path      string  true  "participant web3 address"
// @Success      200  {object}   IfParticipantExistsResp
// @Failure      400  {object}  ErrorMsg
// @Router       /if_participant_exists [get]
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

// UpdataBasicPacticipant godoc
// @Summary     TODO
// @Description TODO
// @Tags		Participants
// @Accept      json
// @Produce     json
// @Param       account body BasicInfoUpdateRequest true "update basic participant info"
// @param		Authorization header string true "Bearer {JWT token}"
// @Success     200  {object}   AuthResp
// @Failure     400  {object}  ErrorMsg
// @Router		/update_basic_info [post]
func (rs *RestSrv) HandleUpdataBasicPacticipant(w http.ResponseWriter, r *http.Request) {
	web3Address, err := rs.getWeb3Address(r)
	if err != nil {
		responError(w, http.StatusBadGateway, err.Error())
		return
	}

	request := new(BasicInfoUpdateRequest)
	if err := rs.getRequest(r.Body, request); err != nil {
		responError(w, http.StatusBadRequest, err.Error())
		return
	}

	participant, err := rs.libSrv.UpdateParticipantNickName(web3Address, request.NickName)
	if err != nil {
		responError(w, http.StatusBadRequest, err.Error())
		return
	}

	// generate a new jwt token for him
	jwt, err := rs.getJWTToken(participant.ID, web3Address, participant.Language, int64(participant.Role))
	if err != nil {
		responError(w, http.StatusBadGateway, "something went wrong, our apologies")
		return
	}

	responJSON(w, http.StatusOK, AuthResp{Token: jwt.Token, Role: participant.Role, NickName: request.NickName})
}

// GetBasicInfo godoc
// @Summary		TODO
// @Description TODO
// @Tags        Participants
// @Accept      json
// @Produce     json
// @Param 		Authorization header string true "Bearer {JWT token}"
// @Success     200  {object}  BasicInfo
// @Failure     400  {object}  ErrorMsg
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
		responError(w, http.StatusBadRequest, err.Error())
		return
	}

	responJSON(w, http.StatusOK, BasicInfo{NickName: participant.NickName, Role: participant.Role})
}
