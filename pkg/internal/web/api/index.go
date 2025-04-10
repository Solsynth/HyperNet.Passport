package api

import (
	"github.com/gofiber/fiber/v2"
)

func MapControllers(app *fiber.App, baseURL string) {
	api := app.Group(baseURL).Name("API")
	{
		api.Get("/well-known/openid-configuration", getOidcConfiguration)
		api.Get("/well-known/jwks", getJwk)

		checkIn := api.Group("/check-in").Name("Daily Check In API")
		{
			checkIn.Get("/", listCheckInRecord)
			checkIn.Get("/today", getTodayCheckIn)
			checkIn.Post("/", doCheckIn)
		}

		notify := api.Group("/notifications").Name("Notifications API")
		{
			notify.Get("/", listNotification)
			notify.Get("/count", getNotificationCount)
			notify.Get("/subscription", getNotifySubscriber)
			notify.Post("/subscription", addNotifySubscriber)
			notify.Delete("/subscription/:deviceId", removeNotifySubscriber)
			notify.Put("/read", markNotificationReadBatch)
			notify.Put("/read/all", markNotificationAllRead)
			notify.Put("/read/:notificationId", markNotificationRead)
		}

		preferences := api.Group("/preferences").Name("Preferences API")
		{
			preferences.Get("/auth", getAuthPreference)
			preferences.Put("/auth", updateAuthPreference)
			preferences.Get("/notifications", getNotificationPreference)
			preferences.Put("/notifications", updateNotificationPreference)
		}

		badges := api.Group("/badges").Name("Badges")
		{
			badges.Get("/me", listUserBadge)
			badges.Post("/:badgeId/active", activeUserBadge)
		}

		reports := api.Group("/reports").Name("Reports API")
		{
			abuse := reports.Group("/abuse").Name("Abuse Reports")
			{
				abuse.Get("/", listAbuseReports)
				abuse.Get("/:id", getAbuseReport)
				abuse.Put("/:id/status", updateAbuseReportStatus)
				abuse.Post("/", createAbuseReport)
			}
		}

		punishments := api.Group("/punishments").Name("Punishments API")
		{
			punishments.Get("/", listUserPunishment)
			punishments.Get("/given", listMadePunishment)
			punishments.Get("/:id", getPunishment)
			punishments.Post("/", createPunishment)
			punishments.Put("/:id", editPunishment)
			punishments.Delete("/:id", deletePunishment)
		}

		api.Get("/users", getUserInBatch)
		api.Get("/users/lookup", lookupAccount)
		api.Get("/users/search", searchAccount)

		me := api.Group("/users/me").Name("Myself Operations")
		{
			me.Get("/avatar", getAvatar)
			me.Get("/banner", getBanner)
			me.Put("/avatar", setAvatar)
			me.Put("/banner", setBanner)

			me.Get("/", getUserinfo)
			me.Get("/oidc", getUserinfoForOidc)
			me.Put("/", editUserinfo)
			me.Put("/language", updateAccountLanguage)
			me.Get("/events", getEvents)
			me.Get("/tickets", getTickets)
			me.Delete("/tickets/:ticketId", deleteTicket)

			me.Post("/confirm", doRegisterConfirm)
			me.Patch("/confirm", reNotifyRegisterConfirm)

			me.Get("/status", getMyselfStatus)
			me.Post("/status", setStatus)
			me.Put("/status", editStatus)
			me.Delete("/status", clearStatus)

			me.Get("/pages", getOwnAccountPage)
			me.Put("/pages", updateAccountPage)

			contacts := me.Group("/contacts").Name("Contacts")
			{
				contacts.Get("/", listContact)
				contacts.Get("/:contactId", getContact)
				contacts.Post("/", createContact)
				contacts.Put("/:contactId", updateContact)
				contacts.Delete("/:contactId", deleteContact)
			}

			factors := me.Group("/factors").Name("Factors")
			{
				factors.Get("/", listFactor)
				factors.Post("/", createFactor)
				factors.Delete("/:factorId", deleteFactor)
			}

			relations := me.Group("/relations").Name("Relations")
			{
				relations.Post("/", makeFriendship)
				relations.Post("/friend", makeFriendship)
				relations.Post("/block", makeBlockship)

				relations.Get("/", listRelationship)
				relations.Get("/:relatedId", getRelationship)
				relations.Put("/:relatedId", editRelationship)
				relations.Delete("/:relatedId", deleteRelationship)

				relations.Post("/:relatedId", makeFriendship)
				relations.Post("/:relatedId/accept", acceptFriend)
				relations.Post("/:relatedId/decline", declineFriend)
			}

			me.Post("/password-reset", requestResetPassword)
			me.Patch("/password-reset", confirmResetPassword)

			me.Post("/deletion", requestDeleteAccount)
			me.Patch("/deletion", confirmDeleteAccount)
		}

		directory := api.Group("/users/:alias").Name("User Directory")
		{
			directory.Get("/", getOtherUserinfo)
			directory.Get("/status", getStatus)
			directory.Get("/page", getAccountPage)

			directory.Get("/check-in", listOtherUserCheckInRecord)
		}

		api.Get("/users", getOtherUserinfoBatch)
		api.Post("/users", doRegister)

		auth := api.Group("/auth").Name("Auth")
		{
			auth.Post("/", doAuthenticate)
			auth.Patch("/", doAuthTicketCheck)
			auth.Post("/token", getToken)

			auth.Get("/tickets/:ticketId", getTicket)

			auth.Get("/factors", getAvailableFactors)
			auth.Post("/factors/:factorId", requestFactorToken)

			auth.Get("/o/authorize", tryAuthorizeThirdClient)
			auth.Post("/o/authorize", authorizeThirdClient)
		}

		realms := api.Group("/realms").Name("Realms API")
		{
			realms.Get("/", listCommunityRealm)
			realms.Get("/me", listOwnedRealm)
			realms.Get("/me/available", listAvailableRealm)
			realms.Get("/:realm", getRealm)
			realms.Get("/:realm/members", listRealmMembers)
			realms.Get("/:realm/members/me", getMyRealmMember)
			realms.Post("/", createRealm)
			realms.Put("/:realmId", editRealm)
			realms.Delete("/:realmId", deleteRealm)
			realms.Post("/:realm/members", addRealmMember)
			realms.Delete("/:realm/members/:memberId", removeRealmMember)
			realms.Delete("/:realm/me", leaveRealm)
		}

		programs := api.Group("/programs").Name("Programs API")
		{
			programs.Get("/", listProgram)
			programs.Get("/members", listProgramMembership)
			programs.Get("/:programId", getProgram)
			programs.Post("/:programId", joinProgram)
			programs.Delete("/:programId", leaveProgram)
		}

		developers := api.Group("/dev").Name("Developers API")
		{
			developers.Post("/notify/:user", notifyUser)
			developers.Post("/notify/all", notifyAllUser)

			bots := developers.Group("/bots").Name("Bots")
			{
				bots.Get("/", listBots)
				bots.Post("/", createBot)
				bots.Delete("/:botId", deleteBot)

				keys := bots.Group("/:botId/keys").Name("Bots' Keys")
				{
					keys.Get("/", listBotKeys)
					keys.Post("/", createBotKey)
					keys.Post("/:id/roll", rollBotKey)
					keys.Put("/:id", editBotKey)
					keys.Delete("/:id", revokeBotKey)
				}
			}

			keys := developers.Group("/keys").Name("Own Bots' Keys")
			{
				keys.Get("/", listBotKeys)
				keys.Get("/:id", getBotKey)
				keys.Post("/", createBotKey)
				keys.Post("/:id/roll", rollBotKey)
				keys.Put("/:id", editBotKey)
				keys.Delete("/:id", revokeBotKey)
			}
		}

		api.Post("/permissions/check", checkPermission)
		api.Post("/permissions/check/:userId", checkUserPermission)

		api.All("/*", func(c *fiber.Ctx) error {
			return fiber.ErrNotFound
		})
	}
}
