package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// tangkap parameter di handler
// handler ke service
// service yang menentukan repo mana yang di panggil
// repo : getAll, getByuserID
// db

type campaignHandler struct{
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

// api/v1/campaigns
func (h *campaignHandler) GetCampaigns(c *gin.Context){
	userID, _ := strconv.Atoi(c.Query("user_id"))
	// strconv.atoi itu untuk convert dari string ke integer

	campaign, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Error to getcampaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaign)
	c.JSON(http.StatusOK, response)
	
}