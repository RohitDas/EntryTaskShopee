package mapper

import (
	"bufio"
	"encoding/json"
	"fmt"
	ads "git.garena.com/shopee-server/shopee_protobuf/beeshop_ads.pb"
	"io"
	"log"
)

func CustomMap(source io.Reader, destination io.Writer) {
	input := bufio.NewScanner(source)
	for input.Scan() {
		line := input.Text()
		var track ads.Tracking
		if err := json.Unmarshal([]byte(line), &track); err != nil {
			log.Fatal("Failed to unmarshal", err.Error())
		} else {
			userId := track.GetUserid()
			//sessionId := track.GetSessionid()
			ts := track.GetTimestamp()
			if track.GetOperation() == int32(ads.TrackingOperationType_CLICK) {
				for _, item := range track.GetItems() {
					shopId := item.GetShopid()
					itemid := item.GetItemid()
					//fmt.Printf("%d\t%d,%d,%d\n", userId, shopId, itemid, ts)
					fmt.Fprint(destination, fmt.Sprintf("%d\t%d,%d,%d\n", userId, shopId, itemid, ts))
				}
			}
		}
	}

}
