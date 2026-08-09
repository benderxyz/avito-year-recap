export type DemoUser = {
  id: string;
  name: string;
  externalId: string;
};

export const DEMO_USERS: DemoUser[] = [
  { id: '42', name: 'Alex', externalId: 'avito-42' },
  { id: '43', name: 'Maria', externalId: 'avito-43' },
  { id: '44', name: 'Ivan', externalId: 'avito-44' },
  { id: '45', name: 'Olga', externalId: 'avito-45' },
  { id: '46', name: 'Dmitry', externalId: 'avito-46' },
  { id: '47', name: 'Elena', externalId: 'avito-47' },
  { id: '48', name: 'Sergey', externalId: 'avito-48' },
  { id: '49', name: 'Anna', externalId: 'avito-49' },
  { id: '50', name: 'Pavel', externalId: 'avito-50' },
  { id: '51', name: 'Nina', externalId: 'avito-51' },
];

export function getDemoUserById(id: string): DemoUser | undefined {
  return DEMO_USERS.find((user) => user.id === id);
}

export function getDemoUserInitial(name: string): string {
  return name.charAt(0).toUpperCase();
}
