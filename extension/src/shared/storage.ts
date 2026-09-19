import { DEFAULT_PREFERENCES, type Application, type Job, type Preferences, type Profile, type Resume } from './types';
const get = async <T>(key:string, fallback:T):Promise<T> => { const result=await chrome.storage.local.get(key); return (result[key] as T|undefined) ?? fallback; };
const set=<T>(key:string,value:T)=>chrome.storage.local.set({[key]:value});
export const storage = {
  profile: () => get<Profile|null>('profile', null), saveProfile: (v:Profile) => set('profile',v),
  preferences: () => get<Preferences>('preferences', DEFAULT_PREFERENCES), savePreferences: (v:Preferences) => set('preferences',v),
  resumes:()=>get<Resume[]>('resumes',[]), saveResumes:(v:Resume[])=>set('resumes',v),
  applications:()=>get<Application[]>('applications',[]), saveApplications:(v:Application[])=>set('applications',v),
  seenJobs: () => get<Record<string,Job>>('seenJobs', {}),
  rememberJob: async (job:Job) => { const jobs=await storage.seenJobs(); const isNew=!jobs[job.id]; jobs[job.id]=job; await set('seenJobs',jobs); return isNew; },
  upsertApplication:async(app:Application)=>{const all=await storage.applications();const i=all.findIndex(a=>a.id===app.id);if(i>=0)all[i]=app;else all.unshift(app);await storage.saveApplications(all)},
  duplicate:async(job:Job)=>{const all=await storage.applications();const key=(s:string)=>s.toLowerCase().replace(/[^a-z0-9]/g,'');return all.find(a=>a.job.url===job.url||(key(a.job.company)===key(job.company)&&key(a.job.title)===key(job.title)))} ,
  clearPersonalData:()=>chrome.storage.local.clear()
};
