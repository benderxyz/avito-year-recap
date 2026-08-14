import { NumberInput, type NumberInputProps } from '@mantine/core';

type NumberFilterProps = Omit<NumberInputProps, 'value' | 'onChange'> & {
  value: number | null;
  onChange: (value: number | null) => void;
};

export default function NumberFilter({
  value,
  onChange,
  label = 'Number',
  ...props
}: NumberFilterProps) {
  return (
    <NumberInput
      {...props}
      label={label}
      value={value ?? ''}
      onChange={(next) => onChange(typeof next === 'number' ? next : null)}
    />
  );
}
