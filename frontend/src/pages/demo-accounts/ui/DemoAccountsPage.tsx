import { Avatar, Card, SimpleGrid, Stack, Text, Title } from '@mantine/core';
import { Link } from '@tanstack/react-router';
import { DEMO_USERS, getDemoUserInitial } from '@/entities/demo-user';

export function DemoAccountsPage() {
  return (
    <div className="page-shell">
      <main className="page-content">
        <Stack gap="xs">
          <Title order={1}>Demo аккаунты</Title>
          <Text c="dimmed">Выберите тестового пользователя для просмотра каталога Авито</Text>
        </Stack>

        <SimpleGrid cols={{ base: 1, xs: 2, sm: 3, md: 4 }} spacing="md">
          {DEMO_USERS.map((user) => (
            <Link
              key={user.id}
              to="/demo/$id"
              params={{ id: user.id }}
              style={{ textDecoration: 'none' }}
            >
              <Card withBorder padding="lg" radius="md">
                <Stack align="center" gap="sm">
                  <Avatar size="lg" radius="xl" color="avito" variant="light">
                    {getDemoUserInitial(user.name)}
                  </Avatar>
                  <Stack gap={2} align="center">
                    <Text fw={700}>{user.name}</Text>
                    <Text size="sm" c="dimmed">
                      ID: {user.id}
                    </Text>
                  </Stack>
                </Stack>
              </Card>
            </Link>
          ))}
        </SimpleGrid>
      </main>
    </div>
  );
}
