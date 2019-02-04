package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var listenPort = flag.Int("listenPort", 1433, "The port which TDS Proxy listens on")
	//var sqlHost = flag.String("sql server address", "localhost", "The host name or IP of the SQL Server")
	//var sqlPort = flag.Int("sql server port", 1433, "The port that SQL Server is running on")
	flag.Parse()

	// System.Net.IPHostEntry iphe = System.Net.Dns.GetHostEntry(args[1]);
	// BridgeAcceptor b = new BridgeAcceptor(
	//     int.Parse(args[0]),
	//     new System.Net.IPEndPoint(iphe.AddressList[0], int.Parse(args[2]))
	// )
	// b.TDSMessageReceived += new TDSMessageReceivedDelegate(b_TDSMessageReceived)
	// b.TDSPacketReceived += new TDSPacketReceivedDelegate(b_TDSPacketReceived)
	// b.ConnectionAccepted += new ConnectionAcceptedDelegate(b_ConnectionAccepted)
	// b.ConnectionDisconnected += new ConnectionDisconnectedDelegate(b_ConnectionClosed)
	// b.Start();
	// fmt.Println("Press enter to kill this process...")
	// Console.ReadLine()
	// b.Stop()

	// catch ^C and clean up
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Exiting Gracefully...")
		//cleanup()
		os.Exit(1)
	}()

	// Listen for incoming connections.
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *listenPort))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Printf("Listening on :%d\n", *listenPort)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}

}

func max(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	var bBuffer []byte
	bHeader := make([]byte, HeaderSize+1, HeaderSize+1)

	fmt.Printf("HERE!\n")
	for bytesRead, _ := conn.Read(bHeader); bytesRead != 0; {
		header := newTDSHeader(bHeader)
		iMinBufferSize := max(0x1000, header.LengthIncludingHeader()+1)

		fmt.Printf("bytesRead: %d\n", bytesRead)
		if (bBuffer == nil) || (HeaderSize < iMinBufferSize) {
			bBuffer = make([]byte, iMinBufferSize, iMinBufferSize)
		}

		fmt.Printf("Header type: %d\n", header.Type())

		if header.Type() == 23 {
			bBuffer = make([]byte, 0x1000-HeaderSize, 0x1000-HeaderSize)
			conn.Read(bBuffer)
			//io.ReadFull(b, bBuffer)
		} else if header.PayloadSize() > 0 {
			//Console.WriteLine("\t{0:N0} bytes package", header.LengthIncludingHeader);
			bBuffer = make([]byte, header.PayloadSize(), header.PayloadSize())
			conn.Read(bBuffer)
			//io.ReadFull(b, bBuffer)
		}
		tdsPacket := newTDSPacket2(bHeader, bBuffer, header.PayloadSize())
		fmt.Println(tdsPacket.String())
	}

	// Close the connection when you're done with it.
	conn.Close()
}

func formatDateTime() string {
	return time.Now().Format(time.StampMicro)
}

// func ConnectionClosed(bc BridgedConnection, ct ConnectionType) {
// 	fmt.Println(FormatDateTime() + "|Connection " + ct + " closed (" + bc.SocketCouple + ")")
// }

// func ConnectionAccepted(sAccepted Socket) {
// 	fmt.Println(FormatDateTime() + "|New connection from " + sAccepted.RemoteEndPoint)
// }

// func TDSPacketReceived(bc BridgedConnection, packet TDSPacket) {
// 	fmt.Println(FormatDateTime() + "|" + packet)
// }

// func TDSMessageReceived(bc BridgedConnection, msg TDSMessage) {
// 	fmt.Println(FormatDateTime() + "|" + msg)
// 	if b, ok := msg.(SQLBatchMessage); ok {
// 		fmt.Print("\tSQLBatch message ")
// 		sstrBatchText := b.GetBatchText()
// 		fmt.Print("({0:N0} chars worth of {1:N0} bytes of data)[", len(strBatchText), len(strBatchText))
// 		fmt.Print(strBatchText)
// 		fmt.Println("]")
// 	} else if rpc, ok := msg.(RPCRequestMessage); ok {
// 		defer func() {
// 			if err := recover(); err != nil {
// 				fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
// 				//os.Exit(1)
// 			}
// 		}()
// 		bPayload := rpc.AssemblePayload()
// 	}
// }
