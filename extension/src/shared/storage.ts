import { DEFAULT_PREFERENCES, type Job, type Preferences, type Profile } from './types';
const get = async <T>(key:string, fallback:T):Promise<T> => (await chrome.storage.local.get(key))[key] ?? fallback;
export const storage = {
  profile: () => get<Profile|null>('profile', null),
  saveProfile: (value:Profile) => chrome.storage.local.set({profile:value}),
  preferences: () => get<Preferences>('preferences', DEFAULT_PREFERENCES),
  savePreferences: (value:Preferences) => chrome.storage.local.set({preferences:value}),
  seenJobs: () => get<Record<string,Job>>('seenJobs', {}),
  rememberJob: async (job:Job) => { const jobs=await storage.seenJobs(); jobs[job.id]=job; await chrome.storage.local.set({seenJobs:jobs}); }
};
