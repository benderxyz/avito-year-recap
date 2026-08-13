import { createFileRoute } from '@tanstack/react-router'
import BadgePage from '@/pages/BadgePage'

export const Route = createFileRoute('/badges/$id')({
  component: BadgePage,
})
