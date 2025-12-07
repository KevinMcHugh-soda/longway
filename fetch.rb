require 'json'
require 'csv'

def get_page(page_number)
  curl = <<-EOS
curl 'https://api.enchor.us/search/advanced' \
  -H 'Accept: application/json, text/plain, */*' \
  -H 'Accept-Language: en-US,en;q=0.9' \
  -H 'Connection: keep-alive' \
  -H 'Content-Type: application/json' \
  -H 'Origin: https://www.enchor.us' \
  -H 'Referer: https://www.enchor.us/' \
  -H 'Sec-Fetch-Dest: empty' \
  -H 'Sec-Fetch-Mode: cors' \
  -H 'Sec-Fetch-Site: same-site' \
  -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36' \
  -H 'sec-ch-ua: "Chromium";v="142", "Google Chrome";v="142", "Not_A Brand";v="99"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "macOS"' \
  --data-raw '{"instrument":null,"difficulty":"expert","drumType":null,"drumsReviewed":false,"sort":null,"source":"website","name":{"value":"","exact":false,"exclude":false},"artist":{"value":"","exact":false,"exclude":false},"album":{"value":"","exact":false,"exclude":false},"genre":{"value":"","exact":false,"exclude":false},"year":{"value":"","exact":false,"exclude":false},"charter":{"value":"Harmonix","exact":true,"exclude":false},"minLength":null,"maxLength":null,"minIntensity":null,"maxIntensity":null,"minAverageNPS":null,"maxAverageNPS":null,"minMaxNPS":null,"maxMaxNPS":null,"minYear":null,"maxYear":null,"modifiedAfter":"","hash":"","trackHash":"","hasSoloSections":null,"hasForcedNotes":null,"hasOpenNotes":null,"hasTapNotes":null,"hasLyrics":null,"hasVocals":null,"hasRollLanes":null,"has2xKick":null,"hasIssues":null,"hasVideoBackground":null,"modchart":null,"page":#{page_number}}'

EOS
  JSON.parse(`#{curl}`)['data'] || []
end

# Accumulate this into downloaded_songs.csv
rows = []
page_number = 1

loop do
  page = get_page(page_number)
  break if page.empty?
  puts "fetching page #{page_number} (#{page.length} songs)"
  page.each do |song|
    rows << {
      id: song['id'],
      title: song['name'],
      artist: song['artist'],
      album: song['album'],
      genre: song['genre'],
      difficulty: song['difficulty'],
      length: song['length'],
      year: song['year'],
    }
  end
  page_number += 1
end

CSV.open('downloaded_songs.csv', 'w') do |csv|
  csv << %w[id title artist album genre difficulty length year]
  rows.each do |row|
    csv << [row[:id], row[:title], row[:artist], row[:album], row[:genre], row[:difficulty], row[:length], row[:year]]
  end
end

puts "Wrote #{rows.length} songs to downloaded_songs.csv"
