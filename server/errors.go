package server

// Rejection reasons: code -> plain-language explanation, per language.
// The portal returns cryptic codes; we map each to why it failed and how to fix it.

type Reason struct {
	Code  string `json:"code"`
	Title string `json:"title"`
	Why   string `json:"why"`
	Fix   string `json:"fix"`
}

var reasons = map[string]map[string]Reason{
	"en": {
		"NAME_MISMATCH":       {Title: "Name does not match Aadhaar", Why: "The name you entered differs from the name on your Aadhaar record (spelling, order or initials).", Fix: "Type your name exactly as printed on your Aadhaar card."},
		"DOB_MISMATCH":        {Title: "Date of birth does not match", Why: "The date of birth differs from the one in your Aadhaar record.", Fix: "Check the date on your Aadhaar card and enter it as DD-MM-YYYY."},
		"MOBILE_NOT_LINKED":   {Title: "Mobile not linked to Aadhaar", Why: "OTP is sent only to the mobile number linked with Aadhaar, and this number is not linked.", Fix: "Use the linked number, or link your mobile at an Aadhaar centre first."},
		"OTP_EXPIRED":         {Title: "OTP expired", Why: "The OTP is valid for a short time and that time has passed.", Fix: "Request a new OTP and enter it within 10 minutes."},
		"OTP_INVALID":         {Title: "Wrong OTP", Why: "The OTP entered does not match the one we sent.", Fix: "Re-enter the latest OTP. Older OTPs stop working when a new one is sent."},
		"PINCODE_INVALID":     {Title: "PIN code is not valid", Why: "A PIN code must be exactly 6 digits and not start with 0.", Fix: "Check the PIN code on your address proof."},
		"PINCODE_MISMATCH":    {Title: "PIN code does not match the location", Why: "The PIN code you entered belongs to a different area than the detected location.", Fix: "Use the detected address, or correct the PIN code."},
		"ADDRESS_INCOMPLETE":  {Title: "Address is incomplete", Why: "Approval needs at least a village/town, district and state that can be verified.", Fix: "Use 'Detect my location' to fill the address automatically."},
		"LOCATION_UNVERIFIED": {Title: "Location could not be verified", Why: "The address could not be matched to a real place, so officials cannot approve it.", Fix: "Use 'Detect my location' instead of typing the address."},
		"DUPLICATE":           {Title: "Application already exists", Why: "An application for this person is already submitted or approved.", Fix: "Check the status of your existing application instead of applying again."},
		"DOCUMENT_UNREADABLE": {Title: "Document is not readable", Why: "The uploaded document is blurred, cropped, or too dark.", Fix: "Upload a clear, well-lit photo showing all four corners."},
		"NOT_ELIGIBLE":        {Title: "Not found in eligibility list", Why: "Your household was not found in the beneficiary database (SECC) for this state.", Fix: "Visit the nearest Ayushman Mitra at a hospital or CSC with your ration card."},
		"ABHA_NOT_FOUND":      {Title: "ABHA number not found", Why: "No ABHA account exists for the number entered.", Fix: "Check the 14-digit number, or create a new ABHA with your Aadhaar."},
		"GATEWAY_UNAVAILABLE": {Title: "Official server is not responding", Why: "The government server is busy or down. Your details are fine.", Fix: "Wait a few minutes and try again."},
		"UNKNOWN":             {Title: "Something went wrong", Why: "The server sent an error we do not recognise.", Fix: "Try again. If it keeps failing, note the code and ask at your nearest Ayushman Mitra."},
	},
	"hi": {
		"NAME_MISMATCH":       {Title: "नाम आधार से मेल नहीं खाता", Why: "आपका लिखा नाम आधार रिकॉर्ड के नाम से अलग है (वर्तनी, क्रम या इनिशियल)।", Fix: "नाम बिल्कुल वैसे लिखें जैसा आधार कार्ड पर छपा है।"},
		"DOB_MISMATCH":        {Title: "जन्मतिथि मेल नहीं खाती", Why: "जन्मतिथि आधार रिकॉर्ड से अलग है।", Fix: "आधार कार्ड की तारीख देखकर DD-MM-YYYY में भरें।"},
		"MOBILE_NOT_LINKED":   {Title: "मोबाइल आधार से जुड़ा नहीं है", Why: "OTP केवल आधार से जुड़े मोबाइल पर जाता है, और यह नंबर जुड़ा नहीं है।", Fix: "जुड़ा हुआ नंबर इस्तेमाल करें, या पहले आधार केंद्र पर मोबाइल जुड़वाएँ।"},
		"OTP_EXPIRED":         {Title: "OTP की समय-सीमा समाप्त", Why: "OTP कुछ ही समय के लिए मान्य होता है और वह समय बीत गया।", Fix: "नया OTP मँगवाएँ और 10 मिनट में भरें।"},
		"OTP_INVALID":         {Title: "गलत OTP", Why: "भरा गया OTP भेजे गए OTP से मेल नहीं खाता।", Fix: "सबसे नया OTP दोबारा भरें। नया OTP आने पर पुराना बेकार हो जाता है।"},
		"PINCODE_INVALID":     {Title: "पिन कोड सही नहीं है", Why: "पिन कोड ठीक 6 अंकों का होना चाहिए और 0 से शुरू नहीं होता।", Fix: "अपने पते के प्रमाण पर पिन कोड जाँचें।"},
		"PINCODE_MISMATCH":    {Title: "पिन कोड स्थान से मेल नहीं खाता", Why: "भरा गया पिन कोड पहचाने गए स्थान के क्षेत्र का नहीं है।", Fix: "पहचाना गया पता इस्तेमाल करें, या पिन कोड सुधारें।"},
		"ADDRESS_INCOMPLETE":  {Title: "पता अधूरा है", Why: "मंज़ूरी के लिए कम से कम गाँव/शहर, ज़िला और राज्य चाहिए।", Fix: "'मेरी लोकेशन पहचानें' दबाकर पता अपने आप भरें।"},
		"LOCATION_UNVERIFIED": {Title: "स्थान सत्यापित नहीं हो सका", Why: "पता किसी असली स्थान से मिलाया नहीं जा सका, इसलिए अधिकारी मंज़ूर नहीं कर सकते।", Fix: "पता हाथ से लिखने के बजाय 'मेरी लोकेशन पहचानें' इस्तेमाल करें।"},
		"DUPLICATE":           {Title: "आवेदन पहले से मौजूद है", Why: "इस व्यक्ति का आवेदन पहले ही जमा या मंज़ूर हो चुका है।", Fix: "दोबारा आवेदन करने के बजाय मौजूदा आवेदन की स्थिति देखें।"},
		"DOCUMENT_UNREADABLE": {Title: "दस्तावेज़ पढ़ा नहीं जा सका", Why: "अपलोड किया दस्तावेज़ धुंधला, कटा हुआ या बहुत अँधेरा है।", Fix: "चारों कोने दिखते हुए साफ़, रोशनी वाली फ़ोटो अपलोड करें।"},
		"NOT_ELIGIBLE":        {Title: "पात्रता सूची में नाम नहीं मिला", Why: "आपके राज्य की लाभार्थी सूची (SECC) में आपका परिवार नहीं मिला।", Fix: "राशन कार्ड लेकर नज़दीकी अस्पताल या CSC में आयुष्मान मित्र से मिलें।"},
		"ABHA_NOT_FOUND":      {Title: "ABHA नंबर नहीं मिला", Why: "दिए गए नंबर का कोई ABHA खाता नहीं है।", Fix: "14 अंकों का नंबर जाँचें, या आधार से नया ABHA बनाएँ।"},
		"GATEWAY_UNAVAILABLE": {Title: "सरकारी सर्वर जवाब नहीं दे रहा", Why: "सरकारी सर्वर व्यस्त या बंद है। आपकी जानकारी सही है।", Fix: "कुछ मिनट बाद फिर कोशिश करें।"},
		"UNKNOWN":             {Title: "कुछ गड़बड़ हुई", Why: "सर्वर ने ऐसी त्रुटि भेजी जिसे हम पहचानते नहीं।", Fix: "फिर कोशिश करें। बार-बार हो तो कोड नोट करके आयुष्मान मित्र से पूछें।"},
	},
}

// lookup returns the reason for code in lang, falling back to English, then UNKNOWN.
func lookup(lang, code string) Reason {
	for _, l := range []string{lang, "en"} {
		if r, ok := reasons[l][code]; ok {
			r.Code = code
			return r
		}
	}
	r := reasons["en"]["UNKNOWN"]
	r.Code = code
	return r
}

func allReasons(lang string) []Reason {
	out := make([]Reason, 0, len(reasons["en"]))
	for code := range reasons["en"] {
		out = append(out, lookup(lang, code))
	}
	return out
}
