import { type ReactNode } from 'react'
import { MantineProvider, mantineTheme } from '../../shared/ui/theme/mantine-theme'

type AppProvidersProps = {
  children?: ReactNode
}

export function AppProviders({ children }: AppProvidersProps) {
  return (
    <MantineProvider theme={mantineTheme} defaultColorScheme="light">
      {children}
    </MantineProvider>
  )
}
