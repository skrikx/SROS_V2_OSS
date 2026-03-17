package trace

import (
	ctrace "srosv2/contracts/trace"
	"srosv2/internal/shared/ids"
)

func LinkReceipt(event *ctrace.TraceEvent, receiptID ids.ReceiptID) {
	event.ReceiptRef = receiptID
}
