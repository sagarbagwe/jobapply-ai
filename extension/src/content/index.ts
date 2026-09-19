import type { Job, Profile } from '../shared/types';
const text=(el:Element|null)=>el?.textContent?.trim()??'';
const meta=(name:string)=>document.querySelector<HTMLMetaElement>(`meta[name="${name}"],meta[property="${name}"]`)?.content??'';
const hash=(s:string)=>{let h=2166136261;for(const c of s){h^=c.charCodeAt(0);h=Math.imul(h,16777619)}return (h>>>0).toString(36)};
function detectJob():Job|null {
  const body=document.body.innerText;
  const title=text(document.querySelector('h1'))||meta('og:title');
  if(!title || !/(responsibilities|qualifications|requirements|about the role|job description)/i.test(body)) return null;
  const company=meta('og:site_name')||text(document.querySelector('[class*="company" i]'))||location.hostname.replace(/^www\./,'');
  const locationText=text(document.querySelector('[class*="location" i]'))||'Not specified';
  const description=text(document.querySelector('[class*="description" i],article,main')).slice(0,30000)||body.slice(0,30000);
  const skills=['Java','Python','Go','JavaScript','React','Node.js','AWS','GCP','Docker','Kubernetes','LLM','LangGraph','Vertex AI'].filter(s=>new RegExp(`\\b${s.replace('.','\\.')}\\b`,'i').test(description));
  return {id:hash(location.href),company,title,location:locationText,url:location.href,description,source:location.hostname,skills};
}
const aliases:Record<string,keyof Profile>={firstname:'fullName',lastname:'fullName',fullname:'fullName',name:'fullName',email:'email',phone:'phone',mobile:'phone',location:'location',city:'location',currentcompany:'currentCompany',currentctc:'currentCtc',expectedctc:'expectedCtc',noticeperiod:'noticePeriod',workauthorization:'workAuthorization'};
function keyFor(el:HTMLInputElement|HTMLTextAreaElement|HTMLSelectElement){const label=el.labels?.[0]?.innerText??el.getAttribute('aria-label')??el.placeholder??el.name??el.id;const k=label.toLowerCase().replace(/[^a-z]/g,'');return Object.entries(aliases).find(([a])=>k.includes(a))?.[1];}
function scanForm(){return [...document.querySelectorAll<HTMLInputElement|HTMLTextAreaElement|HTMLSelectElement>('input:not([type=hidden]),textarea,select')].map((el,index)=>({index,label:el.labels?.[0]?.innerText??el.getAttribute('aria-label')??el.placeholder??el.name??'',type:el instanceof HTMLInputElement?el.type:el.tagName.toLowerCase(),profileKey:keyFor(el)}));}
chrome.runtime.onMessage.addListener((msg,_sender,send)=>{
  if(msg.type==='DETECT_JOB') send({job:detectJob()});
  if(msg.type==='SCAN_FORM') send({fields:scanForm()});
  if(msg.type==='FILL_FORM') { const profile=msg.profile as Profile; let filled=0; document.querySelectorAll<HTMLInputElement|HTMLTextAreaElement|HTMLSelectElement>('input:not([type=hidden]),textarea,select').forEach(el=>{const key=keyFor(el);if(!key||!profile[key]||el.type==='file')return; const value=key==='fullName'&&/first/i.test(el.labels?.[0]?.innerText??'')?profile.fullName.split(' ')[0]:key==='fullName'&&/last/i.test(el.labels?.[0]?.innerText??'')?profile.fullName.split(' ').slice(1).join(' '):String(profile[key]); el.focus(); el.value=value; el.dispatchEvent(new Event('input',{bubbles:true}));el.dispatchEvent(new Event('change',{bubbles:true}));filled++;}); send({filled}); }
  return true;
});
