echo
====


A package which enables use to send ICMP echo requests, to a particular ip. Please note that it is a Toy(a.k.a Hack) implementation and as such not ready for production

Import echo client packge into your project.

    import("github.com/sriluyarlagadda/echo")


Create a new ICMP message with a particular identifier and sequence number,

    echoMessage := echo.NewMessage(10, 245)
    

Set the ip address(IPV4) and the data you want to send

    err := echoMessage.Set(ipAddr, []byte(data))


Receive the response

    response, err := echoMessage.Send()

Done!
