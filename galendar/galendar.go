package galendar

import (
	"context"
	"errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"io/ioutil"
	"log"

	"time"
)

func BuildEvent(summary, location, desc string, startTime, endTime time.Time, timeZone string, attendeeEmails []string) (event *calendar.Event) {

	event = &calendar.Event{
		Summary:     summary,  //title of event
		Location:    location, //"800 Howard St., San Francisco, CA 94103",
		Description: desc,     //"A chance to hear more about Google's developer products.",
		Start: &calendar.EventDateTime{
			DateTime: startTime.Format(time.RFC3339),
			TimeZone: timeZone, //"+0800", or Asia/Singapore
		},
		End: &calendar.EventDateTime{
			DateTime: endTime.Format(time.RFC3339),
			TimeZone: timeZone, // "+0800",
		},

		Reminders: &calendar.EventReminders{
			Overrides: []*calendar.EventReminder{&calendar.EventReminder{
				Method:          "popup",
				Minutes:         10,
				ForceSendFields: nil,
				NullFields:      nil,
			},
				{Method: "email", Minutes: 15},
			},
			UseDefault:      false,
			ForceSendFields: []string{"UseDefault"}, //temp solution to fix google bug
			NullFields:      nil,
		},
		//Attendees: []*calendar.EventAttendee{
		//&calendar.EventAttendee{Email: "abc@outlook.com"},
		//	&calendar.EventAttendee{Email:"sbrin@example.com"},
		//},
		//Recurrence: []string{"RRULE:FREQ=DAILY;COUNT=2"},
	}

	if len(attendeeEmails) > 0 {
		event.Attendees = make([]*calendar.EventAttendee, 0, len(attendeeEmails))

		for _, email := range attendeeEmails {
			event.Attendees = append(event.Attendees, &calendar.EventAttendee{Email: email})
		}
	}

	return
}

func GetCalendarIdNameList(googleSrv *calendar.Service) map[string]string {

	list, err := googleSrv.CalendarList.List().Do()
	if err != nil {
		log.Println("GetCalendarNameIdList err", err)
		return nil
	}

	idName := make(map[string]string)

	for _, item := range list.Items {
		idName[item.Id] = item.Summary
	}
	return idName
}

func GetGoogleOauthConfig(clientSecretFilePath string, scopes ...string) (googleConfig *oauth2.Config) {

	b, err := ioutil.ReadFile(clientSecretFilePath)
	if err != nil {
		log.Println("galendar.GetGoogleOauthConfig Unable to read client secret file:", err)
	}

	gconfig, err := google.ConfigFromJSON(b, scopes...)
	if err != nil {
		log.Println("galendar.GetGoogleOauthConfig Unable to parse client secret file to config:", err)
	}

	return gconfig
}

type GoogleAuthToken struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
}

func GetGoogleCalendarService(gconfig *oauth2.Config, googleToken GoogleAuthToken) (srv *calendar.Service, err error) {
	if gconfig == nil {
		err = errors.New("internal error 399429: cannot connect to google. Please contact support.")

	} else {
		err = errors.New("error 399499: cannot connect to google")
	}

	token := &oauth2.Token{
		AccessToken:  googleToken.AccessToken,
		TokenType:    googleToken.TokenType,
		RefreshToken: googleToken.RefreshToken,
		Expiry:       googleToken.Expiry,
	}

	srv, err = calendar.New(gconfig.Client(context.Background(), token))

	return
}
