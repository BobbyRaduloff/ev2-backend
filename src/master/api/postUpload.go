package api

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/master/env"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/rules"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/sql/dtos/requests"

	dbEmails "github.com/CorporateBusinessTechnologies/email-verifier-v2/src/sql/dtos/emails"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (state *State) PostUpload(c *gin.Context) {
	state.uploadMutex.Lock()
	defer state.uploadMutex.Unlock()

	ffFile, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "please upload a file")
		return
	}

	if !strings.HasSuffix(ffFile.Filename, ".csv") {
		c.String(http.StatusBadRequest, "please upload a csv")
		return
	}

	newFileName := uuid.New().String() + ".csv"

	err = os.MkdirAll(env.API_UPLOADS_DIR, os.ModePerm)
	if err != nil {
		c.String(http.StatusInternalServerError, "could not create uploads directory")
		return
	}

	savePath := filepath.Join(env.API_UPLOADS_DIR, newFileName)
	if err := c.SaveUploadedFile(ffFile, savePath); err != nil {
		c.String(http.StatusInternalServerError, "could not save the file")
		return
	}

	file, err := os.Open(savePath)
	if err != nil {
		c.String(http.StatusInternalServerError, "could not read the file")
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		c.String(http.StatusInternalServerError, "could not read the headers")
		return
	}
	headers[0] = strings.Replace(headers[0], "\xef\xbb\xbf", "", 1)
	headers[0] = strings.TrimPrefix(headers[0], "\ufeff")

	// Read the rest of the lines
	records, err := reader.ReadAll()
	if err != nil {
		c.String(http.StatusInternalServerError, "could not read the rows")
		return
	}

	requestObject := requests.RequestDTO{
		Timestamp:        time.Now(),
		OriginalFilename: ffFile.Filename,
		TotalCount:       len(records),
	}

	requestId, err := requestObject.WriteToDB(state.db)
	if err != nil {
		c.String(http.StatusInternalServerError, "could not create request")
		return
	}

	emailDTOs := make([]dbEmails.EmailDTO, 0, len(records))
	emails := make([]string, 0, len(records))
	emailsMap := make(map[string]bool)
	for _, record := range records {
		apollo := utils.ApolloDataAsBestAsPossible(headers, record)

		isValidEmail := rules.IsEmailValid(apollo.Email)
		if !isValidEmail {
			continue
		}

		if _, ok := emailsMap[apollo.Email]; ok {
			continue
		}

		emailsMap[apollo.Email] = true

		emailObject := dbEmails.EmailDTO{
			Email:         apollo.Email,
			FirstName:     apollo.FirstName,
			LastName:      apollo.LastName,
			Title:         apollo.Title,
			State:         apollo.State,
			City:          apollo.City,
			Country:       apollo.Country,
			CompanyName:   apollo.CompanyName,
			Industry:      apollo.Industry,
			LinkedInLink:  apollo.LinkedInLink,
			EmployeeCount: apollo.EmployeeCount,
			Status:        "QUEUED",
			Timestamp:     time.Now(),
			RequestID:     requestId,
		}

		emailDTOs = append(emailDTOs, emailObject)
		emails = append(emails, apollo.Email)
	}

	chunkedDTOs := utils.ChunkArray(emailDTOs, env.TASK_CHUNK_SIZE)
	chunkedEmails := utils.ChunkArray(emails, env.TASK_CHUNK_SIZE)

	for i, chunk := range chunkedDTOs {
		err = dbEmails.BatchWriteNew(state.db, chunk)
		if err != nil {
			utils.Logger.Error("couldn't upload emails", zap.Error(err))
			c.String(http.StatusInternalServerError, "something went wrong.")
			return
		}
		body := strings.Join(chunkedEmails[i], ",")
		task := fmt.Sprintf("%d,%s", requestId, body)
		log.Println(task)
		state.taskChannel <- task
	}

	c.String(http.StatusOK, "ok")
}
