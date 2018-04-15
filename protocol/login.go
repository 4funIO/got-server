package protocol

import (
	"fmt"
	"log"
)

type login struct {
	rsa RSA
}

func NewLogin(rsa RSA) Protocol {
	return &login{rsa}
}

func (l *login) ReceiveMessage(msg *NetworkMessage) error {
	// TODO(oliverkra): hard coded, must be fixed... missing some stuf of parsePacket (connection.cpp:177)
	msg.position = 7

	// skip client OS
	msg.SkipBytes(2)

	// client version
	version := msg.GetUint16()
	log.Printf("> New client connection using version: %d", version)

	/*
	 * Skipped bytes:
	 * 4 bytes: protocolVersion
	 * 12 bytes: dat, spr, pic signatures (4 bytes each)
	 * 1 byte: 0
	 */
	if version >= 971 {
		msg.SkipBytes(17)
	} else {
		msg.SkipBytes(12)
	}

	if version <= ClientVersionMin {
		return ErrDisconnectUser{
			Message: fmt.Sprintf("Only clients with protocol %s allowed!", ClientVersionStr),
			Version: version,
		}
	}

	if !l.rsa.DecryptNetworkMessage(msg) {
		return ErrDisconnectUser{
			Message: "Are u trying to hack me? :(",
			Version: version,
		}
	}

	key := make([]uint32, 4)
	key[0] = msg.GetUint32()
	key[1] = msg.GetUint32()
	key[2] = msg.GetUint32()
	key[3] = msg.GetUint32()

	log.Println(">> KEY[0]: ", key[0])
	log.Println(">> KEY[1]: ", key[1])
	log.Println(">> KEY[2]: ", key[2])
	log.Println(">> KEY[3]: ", key[3])

	if version < ClientVersionMin || version > ClientVersionMax {
		return ErrDisconnectUser{
			Message: fmt.Sprintf("Only clients with protocol %s allowed!", ClientVersionStr),
			Version: version,
		}
	}

	// if (g_game.getGameState() == GAME_STATE_STARTUP) {
	// 	disconnectClient("Gameworld is starting up. Please wait.", version);
	// 	return;
	// }

	// if (g_game.getGameState() == GAME_STATE_MAINTAIN) {
	// 	disconnectClient("Gameworld is under maintenance.\nPlease re-connect in a while.", version);
	// 	return;
	// }

	// BanInfo banInfo;
	// auto connection = getConnection();
	// if (!connection) {
	// 	return;
	// }

	// if (IOBan::isIpBanned(connection->getIP(), banInfo)) {
	// 	if (banInfo.reason.empty()) {
	// 		banInfo.reason = "(none)";
	// 	}

	// 	std::ostringstream ss;
	// 	ss << "Your IP has been banned until " << formatDateShort(banInfo.expiresAt) << " by " << banInfo.bannedBy << ".\n\nReason specified:\n" << banInfo.reason;
	// 	disconnectClient(ss.str(), version);
	// 	return;
	// }

	accountName := msg.GetString()
	if accountName == "" {
		return ErrDisconnectUser{
			Message: "Invalid account name.",
			Version: version,
		}
	}

	password := msg.GetString()
	if password == "" {
		return ErrDisconnectUser{
			Message: "Invalid password.",
			Version: version,
		}
	}

	msg.SkipBytes(msg.length - 128 - msg.position)
	if !l.rsa.DecryptNetworkMessage(msg) {
		return ErrDisconnectUser{
			Message: "Invalid authentication token.",
			Version: version,
		}
	}

	authToken := msg.GetString()

	log.Printf(">> Account Name: '%s'\n", accountName)
	log.Printf(">> Password: '%s'\n", password)
	log.Printf(">> Authentication Token: '%s'\n", authToken)

	// auto thisPtr = std::static_pointer_cast<ProtocolLogin>(shared_from_this());
	// g_dispatcher.addTask(createTask(std::bind(&ProtocolLogin::getCharacterList, thisPtr, accountName, password, authToken, version)));

	return nil
}
