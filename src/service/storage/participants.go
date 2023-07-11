package storage

func (ss *StorageSrv) UpdateParticipantNickName(participantID, newNickName string) error {
	toUpdate := map[string]interface{}{"nick_name": newNickName}
	return ss.psqlDB.Model(Participant{}).Where("id = ?", participantID).
		UpdateColumns(toUpdate).Error
}

func (ss *StorageSrv) UpdateParticipantRole(id string, newRole ParticipantRole) error {
	toUpdate := map[string]interface{}{"role": newRole}
	if err := ss.psqlDB.Model(Participant{}).Where("id = ?", id).UpdateColumns(toUpdate).
		Error; err != nil {
		return err
	}
	return nil
}
