export type Profile = {
  fullName:string; email:string; phone:string; location:string; yearsExperience:number; skills:string[];
  education?:string; degree?:string; graduationYear?:string; companies?:string[]; jobTitles?:string[];
  projects?:string[]; certifications?:string[]; currentCompany?:string; currentCtc?:string; expectedCtc?:string;
  noticePeriod?:string; preferredLocations?:string[]; workAuthorization?:string; linkedinUrl?:string; githubUrl?:string; portfolioUrl?:string;
};
export type Preferences = { roles:string[]; locations:string[]; technologies:string[]; minSalaryLpa:number; employment:string[]; excluded:string[]; minimumMatch:number; notifications:boolean };
export type Resume = { id:string; name:string; uploadedAt:string; profile:Profile; tags:string[] };
export type Job = { id:string; company:string; title:string; location:string; url:string; description:string; source:string; skills:string[]; experience?:string; employmentType?:string; salary?:string; postedAt?:string; applicationUrl?:string };
export type Match = { matchScore:number; category:'HIGH MATCH'|'GOOD MATCH'|'LOW MATCH'; matchedSkills:string[]; missingSkills:string[]; experienceMatch:boolean; locationMatch:boolean; recommendation:'apply'|'review'|'ignore'; reasons:string[] };
export type ApplicationStatus='Discovered'|'Prepared'|'Applied'|'In Review'|'Interview'|'Rejected'|'Offer'|'Withdrawn';
export type Application={id:string;job:Job;resumeId?:string;matchScore:number;status:ApplicationStatus;discoveredAt:string;appliedAt?:string;notes:string;updatedAt:string};
export const DEFAULT_PREFERENCES: Preferences = {roles:['Software Engineer','SDE 1','SDE 2','Backend Engineer','AI Engineer','AI Agent Engineer'],locations:['Bangalore','Hyderabad','Pune','Mumbai','Remote'],technologies:['Java','Python','Go','JavaScript','React','Node.js','AWS','GCP','Docker','Kubernetes','AI','LLM','LangGraph','Vertex AI'],minSalaryLpa:20,employment:['Full-time'],excluded:['Internship','Unpaid','Contract'],minimumMatch:70,notifications:true};
