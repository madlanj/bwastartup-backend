package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
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

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Error to getcampaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
	
}


func (h *campaignHandler) GetCampaign(c *gin.Context){
	// api/v1/campaign/id
	// handler : mapping id yang di url ke struct input untuk dimasukkan ke service, call formatter
	// service : inputnya struct input => memangkap id di url, manggil repo
	// repo untuk get campaign by ID

	var input campaign.GetCampaignDetailInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := h.service.GetCampaignByID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}


// tangkap parameter dariuser ke  input struct 
// ambil current user dari jwt
// panggil service, parameternya input struct
// panggil repo untuk save data campaign 
func (h *campaignHandler) CreateCampaign(c *gin.Context){
	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors" : errors}

		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", errorMessage)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser

	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to create campaign", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) UpdateCampaign(c *gin.Context){
	var inputID campaign.GetCampaignDetailInput
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaign.CreateCampaignInput

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors" : errors}

		response := helper.APIResponse("Failed to update campaign", http.StatusUnprocessableEntity, "error", errorMessage)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	inputData.User = currentUser

	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to update campaign", http.StatusOK, "success", campaign.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
}


// handler 
// tangkap input dan ubah struct input
// save image ke folder
// service (kondisi poin 2 di repo, panggil repo poin1)
//repository :
// 1. create image/save data image ke dalam table campaign image
// 2. ubah is_primaty true ke false
func(h *campaignHandler) UploadImage(c *gin.Context){
	var input campaign.CreateCampaignImageInput

	err := c.ShouldBind(&input)
	// karna form maka menggunakan S
	if err != nil {
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{
			"is_uploaded" : false,
		}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "failed", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	// currentUser = c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	// images/1-namafile.jpg
	path :=  fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{
			"is_uploaded" : false,
		}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "failed", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.service.SaveCampaignImage(input, path)
	if err != nil {
		data := gin.H{
			"is_uploaded" : false,
		}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "failed", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"is_uploaded" : true,
	}
	response := helper.APIResponse("Campaign image successfuly uploaded", http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)
}