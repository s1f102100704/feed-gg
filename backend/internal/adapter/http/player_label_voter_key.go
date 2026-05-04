package httpadapter

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type PlayerLabelVoterKeyGenerator struct {
	salt string
}

func NewPlayerLabelVoterKeyGenerator(salt string) PlayerLabelVoterKeyGenerator {
	return PlayerLabelVoterKeyGenerator{salt: salt}
}

func (g PlayerLabelVoterKeyGenerator) Generate(
	puuid string,
	identity ClientIdentity,
) string {
	source := puuid + "|" + identity.IP + "|" + identity.UserAgent
	mac := hmac.New(sha256.New, []byte(g.salt))
	_, _ = mac.Write([]byte(source))
	return hex.EncodeToString(mac.Sum(nil))
}
