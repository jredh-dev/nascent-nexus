package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jredh-dev/nexus/services/portal/pkg/fees"
	"github.com/jredh-dev/nexus/services/portal/pkg/models"
)

// AdminGiveawayList renders the admin item management page.
func (h *Handler) AdminGiveawayList(w http.ResponseWriter, r *http.Request) {
	items, err := h.giveawayDB.ListItems("")
	if err != nil {
		log.Printf("Admin: error listing items: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	claims, err := h.giveawayDB.ListClaims("")
	if err != nil {
		log.Printf("Admin: error listing claims: %v", err)
		claims = nil
	}

	h.renderTemplate(w, "admin_giveaway.html", map[string]interface{}{
		"Title":    "Manage Giveaways",
		"Year":     time.Now().Year(),
		"LoggedIn": true,
		"Items":    items,
		"Claims":   claims,
	})
}

// AdminGiveawayNew renders the new item form.
func (h *Handler) AdminGiveawayNew(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "admin_giveaway_edit.html", map[string]interface{}{
		"Title":    "New Item",
		"Year":     time.Now().Year(),
		"LoggedIn": true,
		"IsNew":    true,
	})
}

// AdminGiveawayEdit renders the edit form for an existing item.
func (h *Handler) AdminGiveawayEdit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	item, err := h.giveawayDB.GetItem(id)
	if err != nil || item == nil {
		http.NotFound(w, r)
		return
	}

	fee := fees.CalculateDeliveryDefault(item.DistMiles, item.DriveMinutes)

	h.renderTemplate(w, "admin_giveaway_edit.html", map[string]interface{}{
		"Title":    "Edit: " + item.Title,
		"Year":     time.Now().Year(),
		"LoggedIn": true,
		"IsNew":    false,
		"Item":     item,
		"Fee":      fee,
	})
}

// AdminGiveawaySave handles create/update form submissions for items.
func (h *Handler) AdminGiveawaySave(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	title := strings.TrimSpace(r.FormValue("title"))
	description := strings.TrimSpace(r.FormValue("description"))
	imageURL := strings.TrimSpace(r.FormValue("image_url"))
	condition := models.ItemCondition(r.FormValue("condition"))
	status := models.ItemStatus(r.FormValue("status"))
	milesStr := r.FormValue("dist_miles")
	minutesStr := r.FormValue("drive_minutes")

	if title == "" {
		h.renderTemplate(w, "admin_giveaway_edit.html", map[string]interface{}{
			"Title":    "Edit Item",
			"Year":     time.Now().Year(),
			"LoggedIn": true,
			"IsNew":    id == "",
			"Error":    "Title is required.",
			"Item": &models.Item{
				ID: id, Title: title, Description: description,
				ImageURL: imageURL, Condition: condition, Status: status,
			},
		})
		return
	}

	miles, _ := strconv.ParseFloat(milesStr, 64)
	minutes, _ := strconv.Atoi(minutesStr)

	now := time.Now()

	if id == "" {
		// Create new item.
		item := &models.Item{
			ID:           generateID(),
			Title:        title,
			Description:  description,
			ImageURL:     imageURL,
			Condition:    condition,
			Status:       models.ItemStatusAvailable,
			DistMiles:    miles,
			DriveMinutes: minutes,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if err := h.giveawayDB.CreateItem(item); err != nil {
			log.Printf("Admin: error creating item: %v", err)
			http.Error(w, "Failed to create item", http.StatusInternalServerError)
			return
		}
	} else {
		// Update existing item.
		existing, err := h.giveawayDB.GetItem(id)
		if err != nil || existing == nil {
			http.NotFound(w, r)
			return
		}
		existing.Title = title
		existing.Description = description
		existing.ImageURL = imageURL
		existing.Condition = condition
		existing.Status = status
		existing.DistMiles = miles
		existing.DriveMinutes = minutes

		if err := h.giveawayDB.UpdateItem(existing); err != nil {
			log.Printf("Admin: error updating item %s: %v", id, err)
			http.Error(w, "Failed to update item", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/admin/giveaway", http.StatusSeeOther)
}

// AdminGiveawayDelete deletes an item.
func (h *Handler) AdminGiveawayDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.giveawayDB.DeleteItem(id); err != nil {
		log.Printf("Admin: error deleting item %s: %v", id, err)
	}
	http.Redirect(w, r, "/admin/giveaway", http.StatusSeeOther)
}

// AdminClaimUpdate updates a claim's status (confirm, deliver, cancel).
func (h *Handler) AdminClaimUpdate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	status := models.ClaimStatus(r.FormValue("status"))
	if status == "" {
		http.Error(w, "Status is required", http.StatusBadRequest)
		return
	}

	if err := h.giveawayDB.UpdateClaimStatus(id, status); err != nil {
		log.Printf("Admin: error updating claim %s: %v", id, err)
		http.Error(w, "Failed to update claim", http.StatusInternalServerError)
		return
	}

	// If cancelled, mark the item as available again.
	if status == models.ClaimStatusCancelled {
		claim, err := h.giveawayDB.GetClaim(id)
		if err == nil && claim != nil {
			item, err := h.giveawayDB.GetItem(claim.ItemID)
			if err == nil && item != nil && item.Status == models.ItemStatusClaimed {
				item.Status = models.ItemStatusAvailable
				_ = h.giveawayDB.UpdateItem(item)
			}
		}
	}

	// If delivered, mark item as gone.
	if status == models.ClaimStatusDelivered {
		claim, err := h.giveawayDB.GetClaim(id)
		if err == nil && claim != nil {
			item, err := h.giveawayDB.GetItem(claim.ItemID)
			if err == nil && item != nil {
				item.Status = models.ItemStatusGone
				_ = h.giveawayDB.UpdateItem(item)
			}
		}
	}

	http.Redirect(w, r, "/admin/giveaway", http.StatusSeeOther)
}
