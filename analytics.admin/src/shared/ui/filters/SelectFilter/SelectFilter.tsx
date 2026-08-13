import { Select, type SelectProps } from '@mantine/core';

type SelectFilterProps = Omit<SelectProps, 'value' | 'onChange'> & {
  value: string | null;
  onChange: (value: string | null) => void;
};

export default function SelectFilter({
  value,
  onChange,
  label = 'Filter',
  ...props
}: SelectFilterProps) {
  return <Select {...props} label={label} clearable value={value} onChange={onChange} />;
}
