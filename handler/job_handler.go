package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/user_app/model"
	"github.com/user_app/repository"
)

type JobHandler struct {
	jrepo *repository.JobRepository
}

func NewJobHandler(jrepo *repository.JobRepository) *JobHandler {
	return &JobHandler{
		jrepo: jrepo,
	}
}

// for admin
func (jh *JobHandler) PostJob(c *fiber.Ctx) error {

	var job model.Job
	id := c.Locals("userID")
	userID, ok := id.(uint)
	if !ok {
		// Handle the error case when the type assertion fails
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user ID",
		})
	}
	companyid, err := c.ParamsInt("id")

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"status": "error", "message": "Id is invalid", "data": nil})
	}

	if err := c.BodyParser(&job); err != nil {

		return c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"status": "error", "message": "Error on request", "data": err})

	}
	job.CreatedByID = userID
	job.CompanyID = uint(companyid)

	err = jh.jrepo.CreateJob(job)

	if err != nil {

		return c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"status": "error", "message": "Error on request", "data": err})

	}

	return c.Status(http.StatusCreated).JSON(
		&fiber.Map{"status": "success", "message": "Job created successfully", "data": job})

}

func (jh *JobHandler) GetJobByID(c *fiber.Ctx) error {

	id, err := c.ParamsInt("id")

	if err != nil {

		return c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"status": "error", "message": "Id is invalid", "data": err})

	}

	job, err := jh.jrepo.GetJobByID(uint(id))

	if err != nil {

		return c.Status(http.StatusNotFound).JSON(
			&fiber.Map{"status": "error", "message": "Job not found", "data": err})

	}

	return c.Status(http.StatusOK).JSON(
		&fiber.Map{"status": "success", "message": "Job found", "data": job})

}

func (jh *JobHandler) GetAllJobs(c *fiber.Ctx) error {
	keyword := c.Query("keyword", "")

	jobs, err := jh.jrepo.GetAllJobs(keyword)

	if err != nil {

		return c.Status(http.StatusNotFound).JSON(
			&fiber.Map{"status": "error", "message": "Jobs not found", "data": err})

	}

	return c.Status(http.StatusOK).JSON(
		&fiber.Map{"status": "success", "message": "Jobs found", "data": jobs})

}

func (jh *JobHandler) DeleteJob(c *fiber.Ctx) error {

	var job model.Job

	if err := c.BodyParser(&job); err != nil {

		return c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"status": "error", "message": "Error on request", "data": err})

	}

	err := jh.jrepo.DeleteJob(job)

	if err != nil {

		return c.Status(http.StatusNotFound).JSON(
			&fiber.Map{"status": "error", "message": "Job not found", "data": err})

	}

	return c.Status(http.StatusOK).JSON(
		&fiber.Map{"status": "success", "message": "Job deleted successfully", "data": nil})

}

func (jh *JobHandler) UpdateJob(c *fiber.Ctx) error {

	id, err := c.ParamsInt("id")

	if err != nil {

		return c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"status": "error", "message": "Id is invalid", "data": err})

	}

	var job model.Job

	if err := c.BodyParser(&job); err != nil {

		return c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"status": "error", "message": "Error on request", "data": err})

	}

	err = jh.jrepo.UpdateJob(job, uint(id))

	if err != nil {

		return c.Status(http.StatusNotFound).JSON(
			&fiber.Map{"status": "error", "message": "Job not found", "data": err})

	}

	return c.Status(http.StatusOK).JSON(
		&fiber.Map{"status": "success", "message": "Job updated successfully", "data": job})

}

func (j *JobHandler) GetAdminJobs(c *fiber.Ctx) error {
	adminID := c.Locals("userID").(uint)

	jobs, err := j.jrepo.GetAdminJobs(adminID)

	if err != nil {

		return c.Status(http.StatusNotFound).JSON(
			&fiber.Map{"status": "error", "message": "Jobs not found", "data": err})

	}

	return c.Status(http.StatusOK).JSON(
		&fiber.Map{"status": "success", "message": "Jobs found", "data": jobs})

}
