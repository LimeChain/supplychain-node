package service

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/Limechain/HCS-Integration-Node/app/domain/contract/model"
	proposalRepo "github.com/Limechain/HCS-Integration-Node/app/domain/proposal/repository"
	proposalService "github.com/Limechain/HCS-Integration-Node/app/domain/proposal/service"
	"strings"
)

type WebService struct {
	pr         proposalRepo.ProposalRepository
	ps         *proposalService.ProposalService
	signingKey ed25519.PrivateKey
	peerPubKey ed25519.PublicKey
}

func (s *WebService) Hash(c *model.UnsignedContract) (string, error) {
	proposal, err := s.pr.GetByID(c.ReferencedProposalId)
	if err != nil {
		return "", err
	}

	proposalHash := s.ps.Hash(proposal)

	var sb strings.Builder
	sb.WriteString(c.ContractId)
	sb.WriteRune(',')
	sb.WriteString(c.SupplierId)
	sb.WriteRune(',')
	sb.WriteString(c.BuyerId)
	sb.WriteRune(',')
	sb.WriteString(c.ReferencedProposalId)
	sb.WriteRune(',')
	sb.WriteString(proposalHash)

	return fmt.Sprintf("%x", sha256.Sum256([]byte(sb.String()))), nil

}

func (s *WebService) Sign(c *model.UnsignedContract) (string, error) {
	contractHash, err := s.Hash(c)
	if err != nil {
		return "", err
	}

	signature := ed25519.Sign(s.signingKey, []byte(contractHash))
	signatureStr := hex.EncodeToString(signature)
	return signatureStr, nil
}

func (s *WebService) verify(c *model.UnsignedContract, signature string) (bool, error) {
	contractHash, err := s.Hash(c)
	if err != nil {
		return false, err
	}
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false, err
	}
	return ed25519.Verify(s.peerPubKey, []byte(contractHash), signatureBytes), nil
}

func (s *WebService) VerifyBuyer(c *model.Contract) (bool, error) {
	return s.verify(&c.UnsignedContract, c.BuyerSignature)
}

func (s *WebService) VerifySupplier(c *model.Contract) (bool, error) {
	return s.verify(&c.UnsignedContract, c.SupplierSignature)
}

func New(signingKey ed25519.PrivateKey, pr proposalRepo.ProposalRepository, ps *proposalService.ProposalService, peerPubKey ed25519.PublicKey) *WebService {
	return &WebService{signingKey: signingKey, pr: pr, ps: ps, peerPubKey: peerPubKey}
}
