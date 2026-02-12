# Introduction
This application suppor the following environments:
-  development
-  test
-  production

For now, I'll have a single environment, `development`, with a single file, `.env`.
  
# Location
These environments' configuration files reside on
```bash
internal/config
```

# Files and Order of loading

- Shared by all environments
  - .env, .env.local (not available for the test environment)
- development
  - env.development, .env.development.local
- est
  - env.test, .env.test.local
- production
  - env.production, .env.production.local

The table below depicts the order of loading these files:
<table style="width: 100%; border-collapse: collapse;">
  <thead>
    <tr style="background-color: #010101;">
       <th rowspan="2" style="border: 1px solid #ddd; padding: 8px;">Order</th>
       <th colspan="3" style="border: 1px solid #ddd; padding: 8px;">Environment</th>
      <th rowspan="2" style="border: 1px solid #ddd; padding: 8px;">.gitignore</th>
      <th rowspan="2" style="border: 1px solid #ddd; padding: 8px;">Notes</th>
    </tr>
    <tr style="background-color: #090909;">
      <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">development</th>
      <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">test</th>
      <th style="border: 1px solid #ddd; padding: 8px; text-align: left;">production</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td style="border: 1px solid #ddd; padding: 8px;">1</td>
      <td style="border: 1px solid #ddd; padding: 8px; text-align: left;">.env</td>
      <td style="border: 1px solid #ddd; padding: 8px; text-align: left;">.env</td>
      <td style="border: 1px solid #ddd; padding: 8px; text-align: left">.env</td>
      <td style="border: 1px solid #ddd; padding: 8px; text-align: left;">NO</td>
	  <td style="border: 1px solid #ddd; padding: 8px; text-align: left;">Shared for all environments</td>
    </tr>
    <tr>
      <td style="border: 1px solid #ddd; padding: 8px;">2</td>
      <td style="border: 1px solid #ddd; padding: 8px; text-align: left;">.env.development</td>
      <td style="border: 1px solid #ddd; padding: 8px; text-align: left;">.env.test</td>
      <td style="border: 1px solid #ddd; padding: 8px; text-align: left;">.env.production</td>
      <td style="border: 1px solid #ddd; padding: 8px; text-align: left;">NO</td>
	  <td style="border: 1px solid #ddd; padding: 8px; text-align: left;">Shared environment-specific variables</td>
    </tr>
   <tr>
      <td style="border: 1px solid #ddd; padding: 8px;">3</td>
      <td style="border: 1px solid #ddd; padding: 8px; text-align: left;">.env.local</td>
      <td style="border: 1px solid #ddd; padding: 8px; text-align: left;">N/A</td>
      <td style="border: 1px solid #ddd; padding: 8px; text-align: left;">.env.local</td>
      <td style="border: 1px solid #ddd; padding: 8px; text-align: left;">YES</td>
	  <td style="border: 1px solid #ddd; padding: 8px; text-align: left;">Local overrides</td>
    </tr>
   <tr>
      <td style="border: 1px solid #ddd; padding: 8px;">4</td>
      <td style="border: 1px solid #ddd; padding: 8px; text-align: left;">.env.development.local</td>
      <td style="border: 1px solid #ddd; padding: 8px; text-align: left;">.env.test.local</td>
      <td style="border: 1px solid #ddd; padding: 8px; text-align: left;">..env.production.local</td>
      <td style="border: 1px solid #ddd; padding: 8px; text-align: left;">YES</td>
	  <td style="border: 1px solid #ddd; padding: 8px; text-align: left;">Environment-specific local overrides</td>
    </tr>
  </tbody>
</table>