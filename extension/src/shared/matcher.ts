import type { Job, Match, Preferences, Profile } from './types';
const norm=(s:string)=>s.toLowerCase().replace(/[^a-z0-9+#.]/g,'');
export function matchJob(profile:Profile, prefs:Preferences, job:Job):Match {
  const hay=norm(`${job.title} ${job.description} ${job.skills.join(' ')}`);
  const skills=[...new Set([...profile.skills,...prefs.technologies])];
  const matched=skills.filter(s=>hay.includes(norm(s)));
  const missing=prefs.technologies.filter(s=>!hay.includes(norm(s))).slice(0,8);
  const locationMatch=prefs.locations.some(l=>norm(job.location).includes(norm(l))||norm(l)==='remote'&&hay.includes('remote'));
  const roleMatch=prefs.roles.some(r=>norm(job.title).includes(norm(r))||norm(r).includes(norm(job.title)));
  const excluded=prefs.excluded.some(x=>hay.includes(norm(x)));
  const experienceMatch=!job.experience || !/([4-9]|\d{2})\+?\s*(years|yrs)/i.test(job.experience) || profile.yearsExperience>=Number(job.experience.match(/\d+/)?.[0]??0);
  let score=Math.min(100, matched.length*8+(locationMatch?20:0)+(roleMatch?25:0)+(experienceMatch?15:0));
  if(excluded) score=Math.min(score,25);
  const category=score>=80?'HIGH MATCH':score>=55?'GOOD MATCH':'LOW MATCH';
  return {matchScore:score,category,matchedSkills:matched,missingSkills:missing,experienceMatch,locationMatch,recommendation:excluded||score<45?'ignore':score>=70?'apply':'review',reasons:[roleMatch?'Target role':'Role needs review',locationMatch?'Preferred location':'Location needs review',`${matched.length} matching skills`]};
}
