echo
====


A package which enables use to send ICMP echo requests, to a particular ip. Please note that it is a Toy(a.k.a Hack) implementation and as such not ready for production

Import https://github.com/sriluyarlagadda/echo into your package.
	echoMessage := echo.NewMessage(10, 245)
	err := echoMessage.Set(ipAddr, []byte(data))
	if err != nil {
		fmt.Println("error:", err)
	}

	response, err := echoMessage.Send()
	if err != nil {
		fmt.Println("error:", err)
	}



