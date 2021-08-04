package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/bwmarrin/discordgo"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Don't know if this should really do something
	// We keep it open to satisfy Cloud Run
	// I'd like to be able to use it to clear up old revisions
	// but i do not know how to do that except manually
}

func defaultServer() {
	log.Print("starting server...")
	http.HandleFunc("/", handler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

//var BotID string
func accessSecretVersion() (string, error) {
	project_name := os.Getenv("PROJECT_ID")
	fmt.Println(project_name)
	// name := "projects/my-project/secrets/my-secret/versions/5"
	name := "projects/scratch-project-321714/secrets/bot-token/versions/latest"

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return name, fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return name, fmt.Errorf("failed to access secret version: %v", err)
	}

	secret := string(result.Payload.Data)
	return secret, nil
}

func botSession(Token string) {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	defer dg.Close()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "thanks" {
		s.ChannelMessageSend(m.ChannelID, "np. love you")
	}

	if (m.Content == "Contact on!") && (m.Author.ID == "707723062111371355") {
		s.ChannelMessageSend(m.ChannelID, "hey @shibi#2848")
	} else if m.Content == "Contact on!" {
		s.ChannelMessageSend(m.ChannelID, "hey @rodeo#5783")
	}

	if m.Content == "gm" {
		s.ChannelMessageSend(m.ChannelID, "gm")
	}

	if m.Content == "List on!" {
		s.ChannelMessageSend(m.ChannelID, "<------ TODO LIST for shelny-------->\n    #1 jack of \n <-----end list ------>")
	}
}

func main() {
	Token, _ := accessSecretVersion()

	go defaultServer()

	botSession(Token)

}
