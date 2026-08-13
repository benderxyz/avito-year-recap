import { createFileRoute } from '@tanstack/react-router';
import StoryPage from '@/pages/StoryPage';

export const Route = createFileRoute('/stories/new')({
  component: StoryPage,
});
