package campaign

type Service interface {
	FindCampaigns(userID int) ([]Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) FindCampaigns(userID int) ([]Campaign, error) {
	if userID != 0 {
		campigns, err := s.repository.FindByUserId(userID)
		if err != nil {
			return campigns, err
		}

		return campigns, nil
	}

	campigns, err := s.repository.FindAll()
	if err != nil {
		return campigns, err
	}

	return campigns, nil
}