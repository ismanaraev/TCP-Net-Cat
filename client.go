package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jroimartin/gocui"
)

func drawchat() {
	if len(os.Args) < 3 {
		fmt.Println("[USAGE]: nc $host $port")
		return
	}
	host := os.Args[1]
	port := os.Args[2]
	// Create a new GUI.
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	data := make([]byte, 100)
	user, _ := conn.Read(data)

	var username string
	for {
		fmt.Print(string(data[:user]))
		reader := bufio.NewReader(os.Stdin)

		username, err = reader.ReadString('\n')
		if err != nil {
			log.Print(err)
		}
		if username != "\n" {
			break
		}
	}
	// username = strings.TrimSuffix("username", "\n")

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
		return
	}
	defer g.Close()
	g.Cursor = true

	// Update the views when terminal changes size.
	g.SetManagerFunc(func(g *gocui.Gui) error {
		termwidth, termheight := g.Size()
		_, err := g.SetView("output", 0, 0, termwidth-1, termheight-4)
		if err != nil {
			return err
		}
		_, err = g.SetView("input", 0, termheight-3, termwidth-1, termheight-1)
		if err != nil {
			return err
		}
		return nil
	})

	// Terminal width and height.
	termwidth, termheight := g.Size()

	// Output.
	ov, err := g.SetView("output", 0, 0, termwidth-1, termheight-4)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Failed to create output view:", err)
		return
	}
	ov.Title = " Messages  -  <" + "channel" + "> "

	ov.Autoscroll = true
	ov.Wrap = true

	go func() {
		for {

			data := make([]byte, 1000)
			n, err := conn.Read(data)
			if err != nil {
				log.Fatal(err)
			}

			// fmt.Println(1000)
			// fmt.Print(string(data[:n]))
			fmt.Fprintln(ov, string(data[:n]))
			g.Update(func(g *gocui.Gui) error {
				_, err := g.View("output")
				if err != nil {
					// handle error
				}

				return nil
			})
		}
	}()
	fmt.Fprintf(conn, username)
	// data = make([]byte, 1000)
	// greeting, _ := conn.Read(data)

	// conn.SetReadDeadline(time.Now().Add(time.Second * 1))

	// fmt.Printf(string(data[:greeting]))

	// Send a welcome message.
	// _, err = fmt.Fprintln(ov, string(data[:greeting]))
	// if err != nil {
	// 	log.Println("Failed to print into output view:", err)
	// }
	// _, err = fmt.Fprintln(ov, "<Go-Chat>: Press Ctrl-C to quit.")
	// if err != nil {
	// 	log.Println("Failed to print into output view:", err)
	// }

	// Input.
	iv, err := g.SetView("input", 0, termheight-3, termwidth-1, termheight-1)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Failed to create input view:", err)
		return
	}

	iv.Title = " New Message  -  <" + username + "> "
	iv.FgColor = gocui.ColorWhite
	iv.Editable = true
	err = iv.SetCursor(0, 0)
	if err != nil {
		log.Println("Failed to set cursor:", err)
		return
	}

	// Bind Ctrl-C so the user can quit.
	err = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	})
	if err != nil {
		log.Println("Could not set key binding:", err)
		return
	}

	// Bind enter key to input to send new messages.
	err = g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, iv *gocui.View) error {
		// Read buffer from the beginning.
		fmt.Fprintf(conn, iv.Buffer())

		iv.Rewind()

		// Get output view and print.
		// ov, err := g.View("output")

		if err != nil {
			log.Println("Cannot get output view:", err)
			return err
		}
		// _, err = fmt.Fprintf(ov, "<%s>: %s", username, iv.Buffer())
		// if err != nil {
		// 	log.Println("Cannot print to output view:", err)
		// }

		// Reset input.
		iv.Clear()

		// Reset cursor.
		err = iv.SetCursor(0, 0)
		if err != nil {
			log.Println("Failed to set cursor:", err)
		}
		return err
	})
	if err != nil {
		log.Println("Cannot bind the enter key:", err)
	}

	// Set the focus to input.
	_, err = g.SetCurrentView("input")
	if err != nil {
		log.Println("Cannot set focus to input view:", err)
	}

	// Start the main loop.
	err = g.MainLoop()
	log.Println("Main loop has finished:", err)
}

func main() {
	// Get channel and username.

	// fmt.Print("Enter Desired Username: ")
	// username, err := reader.ReadString('\n')
	// if err != nil {
	// 	log.Println("Could not set username:", err)
	// }
	// Create the GUI.
	drawchat()
}
