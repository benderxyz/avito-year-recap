type AvitoLogoProps = {
  height?: number;
};

export function AvitoLogo({ height = 28 }: AvitoLogoProps) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 120 32"
      height={height}
      aria-label="Авито"
      role="img"
    >
      <title>Авито</title>
      <text
        x="0"
        y="24"
        fill="var(--avito-brand)"
        fontFamily="'Manrope', system-ui, sans-serif"
        fontSize="28"
        fontWeight="800"
        letterSpacing="-0.04em"
      >
        Авито
      </text>
    </svg>
  );
}
