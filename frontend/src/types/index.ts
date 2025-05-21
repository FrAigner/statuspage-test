export interface Service {
  id: string;
  name: string;
  description: string;
  status: 'operational' | 'degraded' | 'outage';
  created_at: string;
  updated_at: string;
}

export interface Incident {
  id: string;
  title: string;
  description: string;
  status: 'investigating' | 'identified' | 'monitoring' | 'resolved';
  impact: 'critical' | 'major' | 'minor';
  service_id: string;
  service?: Service;
  created_at: string;
  updated_at: string;
  resolved_at?: string;
}

export interface IncidentUpdate {
  id: string;
  incident_id: string;
  message: string;
  status: Incident['status'];
  created_at: string;
}
