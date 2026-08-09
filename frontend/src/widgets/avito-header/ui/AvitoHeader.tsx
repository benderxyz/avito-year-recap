import { Avatar, Button, Group, Text, TextInput } from '@mantine/core';
import { useMediaQuery } from '@mantine/hooks';
import { Link } from '@tanstack/react-router';
import { getDemoUserById, getDemoUserInitial } from '@/entities/demo-user';
import { AvitoLogo } from '@/shared/ui';
import styles from './AvitoHeader.module.css';

type AvitoHeaderProps = {
  userId?: string;
  searchValue?: string;
  onSearchChange?: (value: string) => void;
  onSearchSubmit?: () => void;
};

export function AvitoHeader({
  userId,
  searchValue = '',
  onSearchChange,
  onSearchSubmit,
}: AvitoHeaderProps) {
  const user = userId ? getDemoUserById(userId) : undefined;
  const isCompact = useMediaQuery('(max-width: 720px)');

  return (
    <header className={styles.header}>
      <div className={styles.inner}>
        <div className={styles.topRow}>
          <Link to="/demo" className={styles.logoLink}>
            <AvitoLogo height={isCompact ? 22 : 28} />
          </Link>

          <Group gap="sm" className={styles.nav}>
            <Text className={styles.navItem} component="span">
              Избранное
            </Text>
            <Text className={styles.navItem} component="span">
              Сообщения
            </Text>
            <Avatar radius="xl" color="avito" variant="light" size={isCompact ? 'sm' : 'md'}>
              {user ? getDemoUserInitial(user.name) : 'A'}
            </Avatar>
          </Group>
        </div>

        <div className={styles.search}>
          <TextInput
            className={styles.searchInput}
            placeholder="Поиск по объявлениям"
            value={searchValue}
            onChange={(event) => onSearchChange?.(event.currentTarget.value)}
            onKeyDown={(event) => {
              if (event.key === 'Enter') {
                onSearchSubmit?.();
              }
            }}
            variant="unstyled"
            aria-label="Поиск по объявлениям"
          />
          <Button className={styles.searchBtn} onClick={onSearchSubmit}>
            <span className={styles.searchBtnLabel}>Найти</span>
            <span className={styles.searchBtnIcon} aria-hidden="true">
              ↵
            </span>
          </Button>
        </div>
      </div>
    </header>
  );
}
