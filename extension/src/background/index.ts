import { storage } from '../shared/storage';
chrome.runtime.onInstalled.addListener(()=>console.info('JobApply AI installed'));
chrome.runtime.onMessage.addListener((msg,_sender,send)=>{
  if(msg.type==='JOB_SEEN') storage.rememberJob(msg.job).then(()=>send({ok:true}));
  return true;
});
