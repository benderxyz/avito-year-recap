import { TextInput, type TextInputProps } from '@mantine/core';

type SearchFilterProps = Omit<TextInputProps, 'value' | 'onChange'> & {
  value: string;
  onChange: (value: string) => void;
};

export default function SearchFilter({
  value,
  onChange,
  label = 'Search',
  ...props
}: SearchFilterProps) {
  return (
    <TextInput
      {...props}
      label={label}
      value={value}
      onChange={(event) => onChange(event.currentTarget.value)}
    />
  );
}
