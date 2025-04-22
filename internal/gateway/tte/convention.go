package tte

type Convention struct {
	EndDate       string `json:"end_date"`
	GeolocationID string `json:"geolocation_id"`
	ID            string `json:"id"`
	Name          string `json:"name"`
	StartDate     string `json:"start_date"`
	ViewURI       string `json:"view_uri"`
	WebsiteURI    string `json:"website_uri"`

	// AllowAttendeeConversions                  int64    `json:"allow_attendee_conversions"`
	// AllowBadgeBlankLastname                   int64    `json:"allow_badge_blank_lastname"`
	// AllowBadgeEditing                         int64    `json:"allow_badge_editing"`
	// AllowDiscounts                            int64    `json:"allow_discounts"`
	// AllowExhibitorConversions                 int64    `json:"allow_exhibitor_conversions"`
	// AllowGenericTickets                       int64    `json:"allow_generic_tickets"`
	// AllowHostScheduleConflicts                int64    `json:"allow_host_schedule_conflicts"`
	// AllowPermissiveGifting                    int64    `json:"allow_permissive_gifting"`
	// AllowScheduleConflicts                    int64    `json:"allow_schedule_conflicts"`
	// AllowWaitingLists                         int64    `json:"allow_waiting_lists"`
	// ApplyRefundFeeTo                          []string `json:"apply_refund_fee_to"`
	// ApplySalesTaxTo                           []string `json:"apply_sales_tax_to"`
	// BadgeheaderimageID                        *string  `json:"badgeheaderimage_id"`
	// BadgesPerUser                             int64    `json:"badges_per_user"`
	// CanReserveAttendeeSeats                   int64    `json:"can_reserve_attendee_seats"`
	// CanReserveHostSeats                       int64    `json:"can_reserve_host_seats"`
	// Cancelled                                 int64    `json:"cancelled"`
	// ClockType                                 int64    `json:"clock_type"`
	// ContainerAccentColor                      string   `json:"container_accent_color"`
	// ContainerBackgroundColor                  string   `json:"container_background_color"`
	// ContainerTextColor                        string   `json:"container_text_color"`
	// EmailAddress                              string   `json:"email_address"`
	// EndDate                                   string   `json:"end_date"`
	// GeolocationID                             string   `json:"geolocation_id"`
	// GroupID                                   string   `json:"group_id"`
	// ID                                        string   `json:"id"`
	// IsOnline                                  int64    `json:"is_online"`
	// IsSchedulingEnabled                       int64    `json:"is_scheduling_enabled"`
	// IsUsingStripe                             int64    `json:"is_using_stripe"`
	// LibraryID                                 *string  `json:"library_id"`
	// LimitTicketAvailability                   int64    `json:"limit_ticket_availability"`
	// LimitVolunteershiftApplications           int64    `json:"limit_volunteershift_applications"`
	// LinkColor                                 string   `json:"link_color"`
	// MaxBoothsPerExhibitor                     int64    `json:"max_booths_per_exhibitor"`
	// MaxConventionDaysRange                    int64    `json:"max_convention_days_range"`
	// Name                                      string   `json:"name"`
	// PageBackgroundColor                       string   `json:"page_background_color"`
	// PhoneNumber                               *string  `json:"phone_number"`
	// Private                                   int64    `json:"private"`
	// PrototypesEnabled                         int64    `json:"prototypes_enabled"`
	// PurchaserPaysSalesTax                     int64    `json:"purchaser_pays_sales_tax"`
	// RefundDeadlinesBySalesitem                int64    `json:"refund_deadlines_by_salesitem"`
	// RefundFeePercentage                       float64  `json:"refund_fee_percentage"`
	// RestrictedProductsLimitedQuantityPerBadge int64    `json:"restricted_products_limited_quantity_per_badge"`
	// SalesTaxRate                              float64  `json:"sales_tax_rate"`
	// SendExhibitorInfoEmail                    int64    `json:"send_exhibitor_info_email"`
	// ShowAvailableBooths                       int64    `json:"show_available_booths"`
	// ShowAvailableSponsorships                 int64    `json:"show_available_sponsorships"`
	// ShowBadgeSalesCounts                      int64    `json:"show_badge_sales_counts"`
	// ShowSponsorshipSalesCounts                int64    `json:"show_sponsorship_sales_counts"`
	// SkipSkuRelease                            int64    `json:"skip_sku_release"`
	// SlotDuration                              int64    `json:"slot_duration"`
	// SocialmediaimageID                        *string  `json:"socialmediaimage_id"`
	// StartDate                                 string   `json:"start_date"`
	// TapToCollectETickets                      int64    `json:"tap_to_collect_e_tickets"`
	// TicketsPerEventPerBadge                   int64    `json:"tickets_per_event_per_badge"`
	// TwitterHandle                             *string  `json:"twitter_handle"`
	// UpdatesCount                              int64    `json:"updates_count"`
	// URIPart                                   string   `json:"uri_part"`
	// UseDiscord                                int64    `json:"use_discord"`
	// UseETickets                               int64    `json:"use_e_tickets"`
	// VenueID                                   string   `json:"venue_id"`
	// ViewURI                                   string   `json:"view_uri"`
	// WebsiteURI                                string   `json:"website_uri"`
}
