import { NextRequest, NextResponse } from 'next/server';

export async function POST(request: NextRequest) {
  try {
    const { email } = await request.json();

    if (!email || !email.includes('@')) {
      return NextResponse.json(
        { error: 'Valid email is required' },
        { status: 400 }
      );
    }

    const apiKey = process.env.MAILCHIMP_API_KEY;
    const listId = process.env.MAILCHIMP_LIST_ID;

    if (!apiKey || !listId) {
      console.error('Mailchimp credentials not configured');
      return NextResponse.json(
        { error: 'Waitlist service not configured' },
        { status: 500 }
      );
    }

    // Extract datacenter from API key (format: {key}-{dc})
    const datacenter = apiKey.split('-').pop();
    const mailchimpUrl = `https://${datacenter}.api.mailchimp.com/3.0/lists/${listId}/members`;

    const response = await fetch(mailchimpUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${apiKey}`,
      },
      body: JSON.stringify({
        email_address: email,
        status: 'subscribed',
      }),
    });

    const data = await response.json();

    if (!response.ok) {
      // Handle existing subscriber (status 400 with specific error)
      if (response.status === 400 && data.title === 'Member Exists') {
        return NextResponse.json(
          { message: 'You are already on the waitlist!' },
          { status: 200 }
        );
      }
      
      return NextResponse.json(
        { error: data.detail || 'Failed to add to waitlist' },
        { status: response.status }
      );
    }

    return NextResponse.json(
      { message: 'Successfully added to waitlist!' },
      { status: 200 }
    );
  } catch (error) {
    console.error('Waitlist subscription error:', error);
    return NextResponse.json(
      { error: 'An error occurred. Please try again later.' },
      { status: 500 }
    );
  }
}

