import type{Job,Match,Profile,Resume}from'./types';
const BASE='http://localhost:8080/v1';
async function post<T>(path:string,body:unknown):Promise<T>{const r=await fetch(BASE+path,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});if(!r.ok)throw new Error((await r.json().catch(()=>({}))).error??`API ${r.status}`);return r.json()}
export const api={
 parseResume:(name:string,data:string)=>post<{profile:Profile;skills:string[]}>('/resumes/parse',{name,mime_type:'application/pdf',data}),
 match:(profile:Profile,job:Job)=>post<Match>('/match',{profile,job}),
 answer:(profile:Profile,job:Job,question:string)=>post<{answer:string;requires_review:boolean}>('/answers',{profile,job,question}),
 recommend:(job:Job,resumes:Resume[])=>post<{resume_id:string;reason:string}>('/resumes/recommend',{job,resumes:resumes.map(r=>({id:r.id,name:r.name,profile:r.profile,tags:r.tags}))})
};
