import { storage } from '../shared/storage';
chrome.runtime.onInstalled.addListener(()=>console.info('JobApply AI installed'));
chrome.runtime.onMessage.addListener((msg,_sender,send)=>{
  if(msg.type==='JOB_SEEN') (async()=>{const isNew=await storage.rememberJob(msg.job);const prefs=await storage.preferences();if(isNew&&prefs.notifications&&msg.score>=prefs.minimumMatch){await chrome.notifications.create(`job-${msg.job.id}`,{type:'basic',iconUrl:'icon.svg',title:'New Job Match',message:`${msg.job.title} — ${msg.job.company}\nMatch: ${msg.score}%`}).catch(()=>undefined)}send({ok:true,isNew})})();
  return true;
});
