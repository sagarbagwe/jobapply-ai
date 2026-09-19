import type{Application,Job,Match,Profile,Resume}from'./types';
const BASE=(import.meta.env.VITE_API_BASE_URL||'http://localhost:8080/v1').replace(/\/$/,'');
async function token(){return (await chrome.storage.local.get('accessToken')).accessToken as string|undefined}
async function request<T>(path:string,method='GET',body?:unknown):Promise<T>{const access=await token();const r=await fetch(BASE+path,{method,headers:{'Content-Type':'application/json',...(access?{Authorization:`Bearer ${access}`}:{})},body:body===undefined?undefined:JSON.stringify(body)});if(!r.ok)throw new Error((await r.json().catch(()=>({}))).error??`API ${r.status}`);return r.status===204?undefined as T:r.json()}
export const api={
 authenticated:async()=>Boolean(await token()),
 register:async(email:string,password:string)=>{const x=await request<{access_token:string}>('/auth/register','POST',{email,password});await chrome.storage.local.set({accessToken:x.access_token});return x},
 login:async(email:string,password:string)=>{const x=await request<{access_token:string}>('/auth/login','POST',{email,password});await chrome.storage.local.set({accessToken:x.access_token});return x},
 logout:()=>chrome.storage.local.remove('accessToken'),
 parseResume:(name:string,data:string)=>request<{profile:Profile;skills:string[]}>('/resumes/parse','POST',{name,mime_type:'application/pdf',data}),
 match:(profile:Profile,job:Job)=>request<Match>('/match','POST',{profile,job}),
 answer:(profile:Profile,job:Job,question:string)=>request<{answer:string;requires_review:boolean}>('/answers','POST',{profile,job,question}),
 recommend:(job:Job,resumes:Resume[])=>request<{resume_id:string;reason:string}>('/resumes/recommend','POST',{job,resumes:resumes.map(r=>({id:r.id,name:r.name,profile:r.profile,tags:r.tags}))}),
 saveProfile:(profile:Profile)=>request<void>('/profile','PUT',profile),
 track:(a:Application)=>request('/applications','POST',{job:a.job,status:a.status.toLowerCase().replace(' ','_'),match_score:a.matchScore,notes:a.notes}),
 exportData:()=>request<Record<string,unknown>>('/me/export'),
 deleteAccount:async()=>{const access=await token();const r=await fetch(BASE+'/me',{method:'DELETE',headers:{Authorization:`Bearer ${access??''}`,'X-Confirm-Delete':'DELETE MY DATA'}});if(!r.ok)throw new Error('Account deletion failed');await chrome.storage.local.clear()}
};
