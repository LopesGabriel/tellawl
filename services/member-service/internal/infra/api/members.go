package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/domain/models"
	"github.com/lopesgabriel/tellawl/services/member-service/internal/infra/database"
	"go.opentelemetry.io/otel/codes"
)

func (h *apiHandler) HandleGetMemberByID(w http.ResponseWriter, r *http.Request) error {
	ctx, span := h.tracer.Start(r.Context(), "HandleGetMemberByID")
	defer span.End()

	vars := mux.Vars(r)
	memberId := vars["id"]
	member, err := h.repositories.Members.FindByID(ctx, memberId)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			span.SetStatus(codes.Error, "Member not found")
			return NewNotFoundError(ctx, fmt.Sprintf("Member '%s' not found", memberId))
		}

		span.SetStatus(codes.Error, err.Error())
		return err
	}

	result, err := json.Marshal(toMemberAPIResponse(member))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(result)
	span.SetStatus(codes.Ok, "OK")
	return nil
}

func (h *apiHandler) HandleListMembers(w http.ResponseWriter, r *http.Request) error {
	ctx, span := h.tracer.Start(r.Context(), "HandleListMembers")
	defer span.End()

	members := []models.Member{}

	if email := r.URL.Query().Get("email"); email != "" {
		member, err := h.repositories.Members.FindByEmail(ctx, email)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				span.SetStatus(codes.Error, "Member not found")
				return NewNotFoundError(ctx, fmt.Sprintf("Member '%s' not found", email))
			}

			span.SetStatus(codes.Error, err.Error())
			return err
		}

		members = append(members, *member)
	}

	payload := map[string]interface{}{
		"data": toMembersAPIResponse(members),
		"meta": map[string]interface{}{
			"total": len(members),
		},
	}

	result, err := json.Marshal(payload)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(result)
	span.SetStatus(codes.Ok, "OK")
	return nil
}
