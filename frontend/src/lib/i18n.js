import { writable, derived } from 'svelte/store';

export const languages = {
  en: 'English', hi: 'हिन्दी', bn: 'বাংলা', ta: 'தமிழ்', te: 'తెలుగు', mr: 'मराठी'
};

const dict = {
  en: {
    title: 'A-Saarthi', sub: 'Ayushman Bharat, made simple.',
    step1: 'About you', step2: 'Where you live', step3: 'Check and submit',
    name: 'Full name (as on Aadhaar)', dob: 'Date of birth (DD-MM-YYYY)', mobile: 'Mobile linked to Aadhaar',
    detect: 'Detect my location', detecting: 'Detecting…', detectHint: 'No typing needed. We fill the address from your phone’s location.',
    village: 'Village / Town', district: 'District', state: 'State', pincode: 'PIN code',
    next: 'Next', back: 'Back', submit: 'Check my application',
    allGood: 'Everything looks right.', problems: 'Fix these before you submit', why: 'Why', fix: 'How to fix',
    reasons: 'Why applications get rejected', reasonsSub: 'Every reason, in plain words.',
    denied: 'Location permission was denied. Allow it in your browser, or try again.',
    skip: 'Skip to content'
  },
  hi: {
    title: 'ए-सारथी', sub: 'आयुष्मान भारत, आसान तरीके से।',
    step1: 'आपके बारे में', step2: 'आप कहाँ रहते हैं', step3: 'जाँचें और जमा करें',
    name: 'पूरा नाम (आधार के अनुसार)', dob: 'जन्मतिथि (DD-MM-YYYY)', mobile: 'आधार से जुड़ा मोबाइल',
    detect: 'मेरी लोकेशन पहचानें', detecting: 'पहचान रहे हैं…', detectHint: 'टाइप करने की ज़रूरत नहीं। हम फ़ोन की लोकेशन से पता भरते हैं।',
    village: 'गाँव / शहर', district: 'ज़िला', state: 'राज्य', pincode: 'पिन कोड',
    next: 'आगे', back: 'पीछे', submit: 'मेरा आवेदन जाँचें',
    allGood: 'सब कुछ सही है।', problems: 'जमा करने से पहले इन्हें ठीक करें', why: 'क्यों', fix: 'कैसे ठीक करें',
    reasons: 'आवेदन अस्वीकार क्यों होते हैं', reasonsSub: 'हर कारण, सरल भाषा में।',
    denied: 'लोकेशन की अनुमति नहीं मिली। ब्राउज़र में अनुमति दें और फिर कोशिश करें।',
    skip: 'सीधे सामग्री पर जाएँ'
  },
  bn: {
    title: 'এ-সারথি', sub: 'আয়ুষ্মান ভারত, সহজ ভাবে।',
    step1: 'আপনার সম্পর্কে', step2: 'আপনি কোথায় থাকেন', step3: 'যাচাই ও জমা',
    name: 'পূর্ণ নাম (আধার অনুযায়ী)', dob: 'জন্মতারিখ (DD-MM-YYYY)', mobile: 'আধারে যুক্ত মোবাইল',
    detect: 'আমার অবস্থান শনাক্ত করুন', detecting: 'শনাক্ত হচ্ছে…', detectHint: 'টাইপ করতে হবে না। ফোনের অবস্থান থেকে ঠিকানা ভরা হবে।',
    village: 'গ্রাম / শহর', district: 'জেলা', state: 'রাজ্য', pincode: 'পিন কোড',
    next: 'পরবর্তী', back: 'পিছনে', submit: 'আমার আবেদন যাচাই করুন',
    allGood: 'সব ঠিক আছে।', problems: 'জমা দেওয়ার আগে এগুলি ঠিক করুন', why: 'কেন', fix: 'কীভাবে ঠিক করবেন',
    reasons: 'আবেদন কেন বাতিল হয়', reasonsSub: 'প্রতিটি কারণ, সহজ ভাষায়।',
    denied: 'অবস্থানের অনুমতি দেওয়া হয়নি। ব্রাউজারে অনুমতি দিয়ে আবার চেষ্টা করুন।',
    skip: 'বিষয়বস্তুতে যান'
  },
  ta: {
    title: 'ஏ-சாரதி', sub: 'ஆயுஷ்மான் பாரத், எளிமையாக.',
    step1: 'உங்களைப் பற்றி', step2: 'நீங்கள் வசிக்கும் இடம்', step3: 'சரிபார்த்து சமர்ப்பிக்கவும்',
    name: 'முழுப் பெயர் (ஆதார் படி)', dob: 'பிறந்த தேதி (DD-MM-YYYY)', mobile: 'ஆதாருடன் இணைந்த மொபைல்',
    detect: 'என் இருப்பிடத்தைக் கண்டறி', detecting: 'கண்டறிகிறது…', detectHint: 'தட்டச்சு தேவையில்லை. தொலைபேசி இருப்பிடத்திலிருந்து முகவரி நிரப்பப்படும்.',
    village: 'கிராமம் / நகரம்', district: 'மாவட்டம்', state: 'மாநிலம்', pincode: 'அஞ்சல் குறியீடு',
    next: 'அடுத்து', back: 'பின்', submit: 'என் விண்ணப்பத்தைச் சரிபார்',
    allGood: 'எல்லாம் சரியாக உள்ளது.', problems: 'சமர்ப்பிக்கும் முன் இவற்றைச் சரிசெய்யவும்', why: 'ஏன்', fix: 'எப்படிச் சரிசெய்வது',
    reasons: 'விண்ணப்பங்கள் ஏன் நிராகரிக்கப்படுகின்றன', reasonsSub: 'ஒவ்வொரு காரணமும், எளிய மொழியில்.',
    denied: 'இருப்பிட அனுமதி மறுக்கப்பட்டது. உலாவியில் அனுமதித்து மீண்டும் முயலவும்.',
    skip: 'உள்ளடக்கத்திற்குச் செல்'
  },
  te: {
    title: 'ఏ-సారథి', sub: 'ఆయుష్మాన్ భారత్, సులభంగా.',
    step1: 'మీ గురించి', step2: 'మీరు ఉండే చోటు', step3: 'తనిఖీ చేసి సమర్పించండి',
    name: 'పూర్తి పేరు (ఆధార్ ప్రకారం)', dob: 'పుట్టిన తేదీ (DD-MM-YYYY)', mobile: 'ఆధార్‌కు లింక్ అయిన మొబైల్',
    detect: 'నా లొకేషన్‌ను గుర్తించు', detecting: 'గుర్తిస్తోంది…', detectHint: 'టైప్ చేయక్కర్లేదు. ఫోన్ లొకేషన్ నుండి చిరునామా నింపబడుతుంది.',
    village: 'గ్రామం / పట్టణం', district: 'జిల్లా', state: 'రాష్ట్రం', pincode: 'పిన్ కోడ్',
    next: 'తదుపరి', back: 'వెనుకకు', submit: 'నా దరఖాస్తును తనిఖీ చేయి',
    allGood: 'అంతా సరిగ్గా ఉంది.', problems: 'సమర్పించే ముందు వీటిని సరిచేయండి', why: 'ఎందుకు', fix: 'ఎలా సరిచేయాలి',
    reasons: 'దరఖాస్తులు ఎందుకు తిరస్కరించబడతాయి', reasonsSub: 'ప్రతి కారణం, సులభ భాషలో.',
    denied: 'లొకేషన్ అనుమతి నిరాకరించబడింది. బ్రౌజర్‌లో అనుమతించి మళ్లీ ప్రయత్నించండి.',
    skip: 'కంటెంట్‌కు వెళ్ళు'
  },
  mr: {
    title: 'ए-सारथी', sub: 'आयुष्मान भारत, सोप्या पद्धतीने.',
    step1: 'तुमच्याबद्दल', step2: 'तुम्ही कुठे राहता', step3: 'तपासा आणि सबमिट करा',
    name: 'पूर्ण नाव (आधारनुसार)', dob: 'जन्मतारीख (DD-MM-YYYY)', mobile: 'आधारशी जोडलेला मोबाइल',
    detect: 'माझे स्थान ओळखा', detecting: 'ओळखत आहे…', detectHint: 'टाइप करण्याची गरज नाही. फोनच्या स्थानावरून पत्ता भरला जाईल.',
    village: 'गाव / शहर', district: 'जिल्हा', state: 'राज्य', pincode: 'पिन कोड',
    next: 'पुढे', back: 'मागे', submit: 'माझा अर्ज तपासा',
    allGood: 'सर्व काही बरोबर आहे.', problems: 'सबमिट करण्यापूर्वी हे दुरुस्त करा', why: 'का', fix: 'कसे दुरुस्त करावे',
    reasons: 'अर्ज का नाकारले जातात', reasonsSub: 'प्रत्येक कारण, सोप्या भाषेत.',
    denied: 'स्थानाची परवानगी नाकारली. ब्राउझरमध्ये परवानगी द्या आणि पुन्हा प्रयत्न करा.',
    skip: 'मजकुराकडे जा'
  }
};

function initial() {
  try {
    const s = localStorage.getItem('lang');
    if (s && dict[s]) return s;
  } catch {}
  const n = (navigator.language || 'en').slice(0, 2);
  return dict[n] ? n : 'en';
}

export const lang = writable(initial());
lang.subscribe((l) => {
  try { localStorage.setItem('lang', l); } catch {}
  if (typeof document !== 'undefined') document.documentElement.lang = l;
});

export const t = derived(lang, ($l) => (k) => dict[$l]?.[k] ?? dict.en[k] ?? k);
