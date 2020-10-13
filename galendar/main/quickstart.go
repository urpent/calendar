package main

import (
	"encoding/json"
	"fmt"
	myc "github.com/urpent/calendar"

	"github.com/urpent/calendar/galendar"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope) //calendar.CalendarEventsScope,calendar.CalendarReadonlyScope
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	t := time.Date(2020, 7, 15, 1, 30, 0, 0, time.Local).Format(time.RFC3339)
	t2 := time.Now().Add(10 * time.Hour).Format(time.RFC3339)

	calendarId := "c_188bmfgeuf02ei9ok9m3lpgg6nd9q4gbc9nn8rrmclp2sorfdk@resource.calendar.google.com"

	query := srv.Freebusy.Query(&calendar.FreeBusyRequest{
		CalendarExpansionMax: 0,
		GroupExpansionMax:    0,
		Items: []*calendar.FreeBusyRequestItem{&calendar.FreeBusyRequestItem{
			Id:              calendarId,
			ForceSendFields: nil,
			NullFields:      nil,
		}},
		TimeMax:         t2,
		TimeMin:         t,
		TimeZone:        "+0800",
		ForceSendFields: nil,
		NullFields:      nil,
	})

	resp, _ := query.Do()

	g, _ := resp.MarshalJSON()
	log.Println("FreeBusy")

	log.Println(string(g))

	tt := time.Date(2020, 8, 15, 1, 0, 0, 0, time.Local)
	ttEnd := tt.Add(30 * time.Minute)

	log.Println("number of busy: ", len(resp.Calendars[calendarId].Busy))
	log.Println("start:", tt, "   end:", ttEnd)

	for _, item := range resp.Calendars[calendarId].Busy {
		busyStart, _ := time.Parse(time.RFC3339, item.Start)
		busyEnd, _ := time.Parse(time.RFC3339, item.End)

		log.Println("busyStart:", busyStart, "   busyEnd:", busyEnd)
		log.Println("wantedStart:", tt, "   wantedEnd:", ttEnd)

		//log.Println(inTimeSpan(tt, ttEnd, busyStart))
		//log.Println(inTimeSpan(tt, ttEnd, busyEnd))
		log.Println("overlap: ", myc.IsTimeOverlap(tt, ttEnd, busyStart, busyEnd))

		tt = busyEnd
		ttEnd = tt.Add(30 * time.Minute)
	}

	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}

	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {

		for _, item := range events.Items {

			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			fmt.Printf("%v (%v)\n", item.Summary, date)
		}
	}

	start := time.Date(2020, 8, 15, 1, 0, 0, 0, time.Local)
	end := start.Add(30 * time.Minute)

	event := galendar.BuildEvent("Goo1juju",
		"800 Howard St., San Francisco, CA 94103",
		"A chance to hear more about Google's developer products.",
		start,
		end,
		"+0800",
		[]string{"yy_w@outlook.com"})

	//calendarId := "primary"
	event, err = srv.Events.Insert(calendarId, event).SendUpdates("all").Do()
	if err != nil {
		log.Fatalf("Unable to create event. %v\n", err)
	}
	fmt.Printf("Event created: %s\n", event.HtmlLink)
	fmt.Printf("Event status info: %s\n", event.Status)

}

//Can book how many day in advance

func findCalendarFreeSlot(startTime, endTime, betweenStart, betweenEnd time.Time,
	exludedDay []time.Weekday, calendarsBusyPeriods ...[]*calendar.TimePeriod) {

}

//use overlap may be better
//func inTimeSpan(start, end, check time.Time) bool {
//	return (check.Equal(start)|| check.After(start)) && (check.Before(end)|| check.Equal(end))
//}
