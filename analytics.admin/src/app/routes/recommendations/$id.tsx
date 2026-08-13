import { createFileRoute } from '@tanstack/react-router'
import RecommendationPage from '@/pages/RecommendationPage'

export const Route = createFileRoute('/recommendations/$id')({
  component: RecommendationPage,
})
