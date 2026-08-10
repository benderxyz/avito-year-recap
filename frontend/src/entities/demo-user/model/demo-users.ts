export type DemoUser = {
  id: string;
  name: string;
  externalId: string;
};

export function getDemoUserInitial(name: string): string {
  return name.charAt(0).toUpperCase();
}
