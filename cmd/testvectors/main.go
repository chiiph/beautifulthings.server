package main

import (
	"encoding/json"
	"log"
	"os"

	"beautifulthings/account"
)

type accountItem struct {
	Username string
	Password string
	Pk       []byte
	Sk       []byte
}

type encryptionItem struct {
	Input  string
	Output []byte
}

type vectors struct {
	Accounts []accountItem
	Texts    []encryptionItem
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Need output file as an argument")
	}

	accountCases := []struct {
		username string
		password string
	}{
		{"u", "p"},
		{"username1234", "password4321"},
		{"$/()RUFJ)wuwe84349", "dnefu384(/·$873hf7g"},
	}

	var accountResults []accountItem
	var ac *account.Account

	for _, c := range accountCases {
		a, err := account.New(c.username, c.password)
		if err != nil {
			log.Fatalf("Error creating account: %s", err)
		}
		accountResults = append(accountResults, accountItem{
			Username: a.Username,
			Password: c.password,
			Pk:       a.Pk[:],
			Sk:       a.Sk[:],
		})
		ac = a
	}

	encCases := []string{
		"",
		"word",
		"some text a bit longer",
		"sómething with ñ and ü",
		`Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem 
		aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo. 
		Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores 
		eos qui ratione voluptatem sequi nesciunt. Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet, 
		consectetur, adipisci velit, sed quia non numquam eius modi tempora incidunt ut labore et dolore magnam aliquam 
		quaerat voluptatem. Ut enim ad minima veniam, quis nostrum exercitationem ullam corporis suscipit laboriosam, 
		nisi ut aliquid ex ea commodi consequatur? Quis autem vel eum iure reprehenderit qui in ea voluptate velit esse 
		quam nihil molestiae consequatur, vel illum qui dolorem eum fugiat quo voluptas nulla pariatur?`,
	}

	var encItems []encryptionItem

	for _, in := range encCases {
		ct, err := ac.Encrypt(in)
		if err != nil {
			log.Fatalf("Error encrypting text: %s", err)
		}
		encItems = append(encItems, encryptionItem{
			Input:  in,
			Output: ct,
		})
	}

	results := vectors{
		Accounts: accountResults,
		Texts:    encItems,
	}
	b, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling accountResults: %s", err)
	}
	f, err := os.Create(os.Args[1])
	if err != nil {
		log.Fatalf("Error opening file %s: %s", os.Args[1], err)
	}
	defer f.Close()

	f.Write(b)
}
